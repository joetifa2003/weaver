#!/usr/bin/env python3
import psutil
import json
import time
import argparse
import signal
import subprocess
from datetime import datetime

class MemoryMonitor:
    def __init__(self, port, output_file, interval=1, pm2_mode=False, gunicorn_mode=False):
        self.port = port
        self.output_file = output_file
        self.interval = interval
        self.pm2_mode = pm2_mode
        self.gunicorn_mode = gunicorn_mode
        self.monitoring_data = []
        self.start_time = time.time()
        self.running = True
        self.process_objects = {}  # Cache process objects for CPU monitoring
        
        # Set up signal handlers for graceful shutdown
        signal.signal(signal.SIGINT, self.signal_handler)
        signal.signal(signal.SIGTERM, self.signal_handler)
    
    def signal_handler(self, signum, frame):
        """Handle shutdown signals gracefully"""
        print(f"\\nReceived signal {signum}, stopping memory monitoring...")
        self.running = False
    
    def save_data(self):
        """Save memory and CPU data to JSON file"""
        try:
            with open(self.output_file, 'w') as f:
                json.dump(self.monitoring_data, f, indent=2)
            print(f"Memory and CPU data saved to {self.output_file} ({len(self.monitoring_data)} data points)")
        except Exception as e:
            print(f"Error saving data: {e}")
    
    def find_port_processes(self):
        """Find processes using the specified port"""
        pids = []
        try:
            for conn in psutil.net_connections():
                if hasattr(conn, 'laddr') and conn.laddr and conn.laddr.port == self.port:
                    if conn.status == 'LISTEN' and conn.pid:
                        pids.append(conn.pid)
        except (psutil.AccessDenied, psutil.NoSuchProcess):
            pass
        
        if not pids:
            # Try alternative method - look for any process listening on the port
            try:
                import subprocess
                result = subprocess.run(['lsof', '-ti', f':{self.port}'], 
                                      capture_output=True, text=True, timeout=5)
                if result.returncode == 0 and result.stdout.strip():
                    pids = [int(pid) for pid in result.stdout.strip().split('\\n') if pid.strip()]
            except:
                pass
        
        return pids
    
    def find_gunicorn_processes(self):
        """Find Gunicorn master and worker processes by iterating through all processes."""
        gunicorn_pids = set()
        target_port_str = str(self.port)
        
        for proc in psutil.process_iter(['pid', 'name', 'cmdline']):
            try:
                # Check if it's a gunicorn process
                if 'gunicorn' in proc.name().lower() or any('gunicorn' in s.lower() for s in proc.cmdline()):
                    # Check if it's listening on the target port
                    for conn in proc.connections(kind='inet'):
                        if conn.laddr.port == self.port:
                            gunicorn_pids.add(proc.pid)
                            break # Found a connection on the port, no need to check other connections for this process
                    # Also check if the command line explicitly contains the bind address with the port
                    if not gunicorn_pids or proc.pid not in gunicorn_pids:
                        for arg in proc.cmdline():
                            if target_port_str in arg and ('--bind' in arg or '-b' in arg):
                                gunicorn_pids.add(proc.pid)
                                break
            except (psutil.NoSuchProcess, psutil.AccessDenied, psutil.ZombieProcess):
                continue
        
        print(f"DEBUG: Found Gunicorn PIDs: {list(gunicorn_pids)}")
        return list(gunicorn_pids)

    def find_pm2_pids(self, app_name=None):
        try:
            result = subprocess.run(['pm2', 'jlist'], capture_output=True, text=True, timeout=5)
            if result.returncode != 0:
                print("Failed to run pm2 jlist")
                return []
            pm2_list = json.loads(result.stdout)
            pids = []
            for proc in pm2_list:
                if app_name is None or proc['pm2_env']['name'] == app_name:
                    pid = proc.get('pid')
                    if pid and pid > 0:
                        pids.append(pid)
            return pids
        except Exception as e:
            print(f"Error getting pm2 pids: {e}")
            return []

    def monitor_memory(self):
        """Monitor memory and CPU usage based on mode (port-based or pm2)"""
        if self.gunicorn_mode:
            print(f"Monitoring memory and CPU usage for Gunicorn processes on port {self.port}...")
        elif self.pm2_mode:
            print(f"Monitoring memory and CPU usage for all Node.js processes managed by pm2...")
        else:
            print(f"Monitoring memory and CPU usage for processes on port {self.port}...")

        # Get PIDs based on monitoring mode
        if self.gunicorn_mode:
            pids = self.find_gunicorn_processes()
        elif self.pm2_mode:
            pids = self.find_pm2_pids()
        else:
            pids = self.find_port_processes()
        
        try:
            while self.running:
                current_time = time.time()
                elapsed = current_time - self.start_time
                
                total_memory_mb = 0
                total_cpu_percent = 0
                process_count = 0
                
                # Update process objects cache
                current_processes = {}
                for pid in pids:
                    try:
                        if pid in self.process_objects:
                            process = self.process_objects[pid]
                        else:
                            process = psutil.Process(pid)
                            # Call cpu_percent() once to initialize baseline
                            process.cpu_percent()
                        current_processes[pid] = process
                    except (psutil.NoSuchProcess, psutil.AccessDenied):
                        continue
                
                self.process_objects = current_processes
                
                # If this is the first measurement, wait a bit for CPU calculation
                if elapsed < 1:
                    time.sleep(0.1)
                
                for pid, process in self.process_objects.items():
                    try:
                        memory_info = process.memory_info()
                        memory_mb = memory_info.rss / (1024 * 1024)  # Convert to MB
                        cpu_percent = process.cpu_percent()
                        
                        total_memory_mb += memory_mb
                        total_cpu_percent += cpu_percent
                        process_count += 1
                    except (psutil.NoSuchProcess, psutil.AccessDenied):
                        continue
                
                if process_count > 0:
                    data_point = {
                        'timestamp': datetime.now().isoformat(),
                        'elapsed': elapsed,
                        'memory_mb': total_memory_mb,
                        'cpu_percent': total_cpu_percent,
                        'process_count': process_count,
                        'monitoring_mode': 'gunicorn' if self.gunicorn_mode else ('pm2' if self.pm2_mode else 'port')
                    }
                    self.monitoring_data.append(data_point)
                    
                    mode_desc = "Gunicorn processes" if self.gunicorn_mode else ("Node.js processes (pm2)" if self.pm2_mode else f"processes on port {self.port}")
                    print(f"Time: {elapsed:.1f}s, Memory: {total_memory_mb:.2f}MB, CPU: {total_cpu_percent:.2f}%, {mode_desc}: {process_count}")
                elif elapsed > 2:  # Only warn after initial startup period
                    if self.gunicorn_mode:
                        print(f"Time: {elapsed:.1f}s, No Gunicorn processes found on port {self.port}")
                    elif self.pm2_mode:
                        print(f"Time: {elapsed:.1f}s, No Node.js processes found managed by pm2")
                    else:
                        print(f"Time: {elapsed:.1f}s, No processes found on port {self.port}")
                
                # Sleep in small increments to respond quickly to shutdown signals
                sleep_time = 0
                while sleep_time < self.interval and self.running:
                    time.sleep(0.1)
                    sleep_time += 0.1
                
        except Exception as e:
            print(f"Error during monitoring: {e}")
        finally:
            # Always save data before exiting
            self.save_data()

def main():
    parser = argparse.ArgumentParser(description='Monitor memory usage of processes on a specific port or pm2 managed processes')
    parser.add_argument('port', type=int, help='Port number to monitor (ignored in pm2 mode)')
    parser.add_argument('output_file', help='Output JSON file')
    parser.add_argument('--interval', type=float, default=0.1, help='Monitoring interval in seconds (default: 1.0)')
    parser.add_argument('--pm2', action='store_true', help='Monitor all Node.js processes managed by pm2 instead of port-based monitoring')
    parser.add_argument('--gunicorn', action='store_true', help='Monitor Gunicorn master and worker processes')
    
    args = parser.parse_args()
    
    monitor = MemoryMonitor(args.port, args.output_file, args.interval, args.pm2, args.gunicorn)
    monitor.monitor_memory()

if __name__ == "__main__":
    main()
