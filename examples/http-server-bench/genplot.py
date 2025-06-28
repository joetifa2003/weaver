import json
import matplotlib.pyplot as plt
import numpy as np
import argparse
import os
from pathlib import Path
import hashlib

def parse_time_unit(time_str):
    """Convert time string with units to microseconds"""
    if not time_str:
        return 0
  
    time_str = time_str.strip()
    if time_str.endswith('µs'):
        return float(time_str[:-2])
    elif time_str.endswith('ms'):
        return float(time_str[:-2]) * 1000
    elif time_str.endswith('s'):
        return float(time_str[:-1]) * 1000000
    else:
        # Assume microseconds if no unit
        return float(time_str)

def parse_elapsed_time(elapsed_str):
    """Convert elapsed time to seconds"""
    if elapsed_str.endswith('ms'):
        return float(elapsed_str[:-2]) / 1000
    elif elapsed_str.endswith('s'):
        return float(elapsed_str[:-1])
    else:
        return float(elapsed_str)

def get_consistent_color(filename):
    """Generate a consistent color for a filename using hash"""
    # Create a hash of the filename
    hash_object = hashlib.md5(filename.encode())
    hash_hex = hash_object.hexdigest()
  
    # Convert first 6 characters of hash to RGB
    r = int(hash_hex[0:2], 16) / 255.0
    g = int(hash_hex[2:4], 16) / 255.0
    b = int(hash_hex[4:6], 16) / 255.0
  
    return (r, g, b)

def create_color_map(filenames):
    """Create a consistent color mapping for all filenames"""
    color_map = {}
  
    # Define a set of distinct colors for better visibility
    predefined_colors = [
        '#1f77b4',  # blue
        '#ff7f0e',  # orange
        '#2ca02c',  # green
        '#d62728',  # red
        '#9467bd',  # purple
        '#8c564b',  # brown
        '#e377c2',  # pink
        '#7f7f7f',  # gray
        '#bcbd22',  # olive
        '#17becf',  # cyan
        '#aec7e8',  # light blue
        '#ffbb78',  # light orange
        '#98df8a',  # light green
        '#ff9896',  # light red
        '#c5b0d5',  # light purple
    ]
  
    # Sort filenames for consistent ordering
    sorted_filenames = sorted(filenames)
  
    for i, filename in enumerate(sorted_filenames):
        if i < len(predefined_colors):
            color_map[filename] = predefined_colors[i]
        else:
            # Fall back to hash-based color for additional files
            color_map[filename] = get_consistent_color(filename)
  
    return color_map

def load_benchmark_data(file_path):
    """Load and parse benchmark data from a file"""
    data_points = []

    with open(file_path, 'r') as f:
        content = f.read().strip()
    
    # Split by lines and parse each JSON object
    json_objects = []
    current_json = ""
    brace_count = 0

    for line in content.split('\n'):
        line = line.strip()
        if not line:
            continue
        
        current_json += line + '\n'
        brace_count += line.count('{') - line.count('}')
    
        if brace_count == 0 and current_json.strip():
            try:
                json_obj = json.loads(current_json.strip())
                json_objects.append(json_obj)
                current_json = ""
            except json.JSONDecodeError:
                pass

    for obj in json_objects:
        elapsed = parse_elapsed_time(obj['Summary']['Elapsed'])
        rps = obj['Summary']['RPS']
        mean_latency = parse_time_unit(obj['Statistics']['Latency']['Mean'])
        p95_latency = parse_time_unit(obj['Percentiles']['P95'])
        p99_latency = parse_time_unit(obj['Percentiles']['P99'])
    
        data_points.append({
            'elapsed': elapsed,
            'rps': rps,
            'mean_latency': mean_latency / 1000,  # Convert to ms
            'p95_latency': p95_latency / 1000,    # Convert to ms
            'p99_latency': p99_latency / 1000     # Convert to ms
        })

    return data_points

def load_memory_data(mem_file_path):
    """Load and parse memory and CPU data from a memory file"""
    if not os.path.exists(mem_file_path):
        print(f"Warning: Memory file {mem_file_path} not found")
        return []
    
    try:
        with open(mem_file_path, 'r') as f:
            memory_data = json.load(f)
        
        # Convert to consistent format
        data_points = []
        for point in memory_data:
            data_points.append({
                'elapsed': point['elapsed'],
                'memory_mb': point['memory_mb'],
                'cpu_percent': point.get('cpu_percent', 0),  # Default to 0 if not present
                'process_count': point.get('process_count', 1)
            })
        
        print(f"Loaded {len(data_points)} memory/CPU data points from {mem_file_path}")
        return data_points
    except Exception as e:
        print(f"Error loading memory data from {mem_file_path}: {e}")
        return []

def load_directory_data(directory_path):
    """Load both stats and memory data from a directory"""
    stats_file = os.path.join(directory_path, 'stats.json')
    mem_file = os.path.join(directory_path, 'mem.json')
    
    if not os.path.exists(stats_file):
        print(f"Warning: Stats file {stats_file} not found")
        return None, None
    
    print(f"Loading data from {directory_path}...")
    stats_data = load_benchmark_data(stats_file)
    memory_data = load_memory_data(mem_file)
    
    return stats_data, memory_data

def plot_comparisons(dir_data_dict, output_dir='plots'):
    """Create comparison plots for all metrics including memory and CPU"""
    os.makedirs(output_dir, exist_ok=True)
  
    plt.style.use('default')
  
    # Create consistent color mapping for all directories
    color_map = create_color_map(dir_data_dict.keys())
  
    # Plot RPS over time
    plt.figure(figsize=(12, 8))
    for dirname, (stats_data, _) in dir_data_dict.items():
        if stats_data:
            elapsed_times = [point['elapsed'] for point in stats_data]
            rps_values = [point['rps'] for point in stats_data]
            plt.plot(elapsed_times, rps_values, marker='o', label=dirname, 
                    color=color_map[dirname], linewidth=2, markersize=4)
  
    plt.xlabel('Elapsed Time (seconds)')
    plt.ylabel('Requests Per Second (RPS)')
    plt.title('RPS Comparison Over Time')
    plt.legend()
    plt.grid(True, alpha=0.3)
    plt.tight_layout()
    plt.savefig(f'{output_dir}/rps_comparison.svg', bbox_inches='tight')
    plt.close()
  
    # Plot Mean Latency over time
    plt.figure(figsize=(12, 8))
    for dirname, (stats_data, _) in dir_data_dict.items():
        if stats_data:
            elapsed_times = [point['elapsed'] for point in stats_data]
            latencies = [point['mean_latency'] for point in stats_data]
            plt.plot(elapsed_times, latencies, marker='o', label=dirname, 
                    color=color_map[dirname], linewidth=2, markersize=4)
  
    plt.xlabel('Elapsed Time (seconds)')
    plt.ylabel('Mean Latency (ms)')
    plt.title('Mean Latency Comparison Over Time')
    plt.legend()
    plt.grid(True, alpha=0.3)
    plt.tight_layout()
    plt.savefig(f'{output_dir}/mean_latency_comparison.svg', bbox_inches='tight')
    plt.close()
  
    # Plot P95 Latency over time
    plt.figure(figsize=(12, 8))
    for dirname, (stats_data, _) in dir_data_dict.items():
        if stats_data:
            elapsed_times = [point['elapsed'] for point in stats_data]
            p95_latencies = [point['p95_latency'] for point in stats_data]
            plt.plot(elapsed_times, p95_latencies, marker='o', label=dirname, 
                    color=color_map[dirname], linewidth=2, markersize=4)
  
    plt.xlabel('Elapsed Time (seconds)')
    plt.ylabel('P95 Latency (ms)')
    plt.title('P95 Latency Comparison Over Time')
    plt.legend()
    plt.grid(True, alpha=0.3)
    plt.tight_layout()
    plt.savefig(f'{output_dir}/p95_latency_comparison.svg', bbox_inches='tight')
    plt.close()
  
    # Plot P99 Latency over time
    plt.figure(figsize=(12, 8))
    for dirname, (stats_data, _) in dir_data_dict.items():
        if stats_data:
            elapsed_times = [point['elapsed'] for point in stats_data]
            p99_latencies = [point['p99_latency'] for point in stats_data]
            plt.plot(elapsed_times, p99_latencies, marker='o', label=dirname, 
                    color=color_map[dirname], linewidth=2, markersize=4)
  
    plt.xlabel('Elapsed Time (seconds)')
    plt.ylabel('P99 Latency (ms)')
    plt.title('P99 Latency Comparison Over Time')
    plt.legend()
    plt.grid(True, alpha=0.3)
    plt.tight_layout()
    plt.savefig(f'{output_dir}/p99_latency_comparison.svg', bbox_inches='tight')
    plt.close()
    
    # Plot Memory Usage over time
    plt.figure(figsize=(12, 8))
    for dirname, (_, memory_data) in dir_data_dict.items():
        if memory_data:
            elapsed_times = [point['elapsed'] for point in memory_data]
            memory_values = [point['memory_mb'] for point in memory_data]
            plt.plot(elapsed_times, memory_values, marker='o', label=dirname, 
                    color=color_map[dirname], linewidth=2, markersize=4)
  
    plt.xlabel('Elapsed Time (seconds)')
    plt.ylabel('Memory Usage (MB)')
    plt.title('Memory Usage Comparison Over Time')
    plt.legend()
    plt.grid(True, alpha=0.3)
    plt.tight_layout()
    plt.savefig(f'{output_dir}/memory_comparison.svg', bbox_inches='tight')
    plt.close()
    
    # Plot CPU Usage over time
    plt.figure(figsize=(12, 8))
    for dirname, (_, memory_data) in dir_data_dict.items():
        if memory_data:
            elapsed_times = [point['elapsed'] for point in memory_data]
            cpu_values = [point['cpu_percent'] for point in memory_data]
            plt.plot(elapsed_times, cpu_values, marker='o', label=dirname, 
                    color=color_map[dirname], linewidth=2, markersize=4)
  
    plt.xlabel('Elapsed Time (seconds)')
    plt.ylabel('CPU Usage (%)')
    plt.title('CPU Usage Comparison Over Time')
    plt.legend()
    plt.grid(True, alpha=0.3)
    plt.tight_layout()
    plt.savefig(f'{output_dir}/cpu_comparison.svg', bbox_inches='tight')
    plt.close()
    
    # Create summary bar chart
    create_summary_chart(dir_data_dict, output_dir, color_map)

def create_summary_chart(dir_data_dict, output_dir, color_map):
    """Create a summary bar chart with average metrics"""
    languages = []
    avg_rps = []
    avg_latency = []
    avg_p95_latency = []
    avg_memory = []
    colors = []
    
    for dirname, (stats_data, memory_data) in dir_data_dict.items():
        if stats_data:
            languages.append(dirname)
            colors.append(color_map[dirname])
            
            # Calculate averages
            rps_values = [point['rps'] for point in stats_data]
            latency_values = [point['mean_latency'] for point in stats_data]
            p95_latency_values = [point['p95_latency'] for point in stats_data]  # Changed from p99 to p95
            
            avg_rps.append(np.mean(rps_values))
            avg_latency.append(np.mean(latency_values))
            avg_p95_latency.append(np.mean(p95_latency_values))  # Changed from p99 to p95
            
            if memory_data:
                memory_values = [point['memory_mb'] for point in memory_data]
                avg_memory.append(np.mean(memory_values))
            else:
                avg_memory.append(0)
    
    if not languages:
        print("No data available for summary chart")
        return
    
    # Create subplots for summary (2x2 grid)
    fig, ((ax1, ax2), (ax3, ax4)) = plt.subplots(2, 2, figsize=(12, 10))
    
    # Average RPS
    bars1 = ax1.bar(languages, avg_rps, color=colors)
    ax1.set_title('Average RPS')
    ax1.set_ylabel('Requests Per Second')
    ax1.tick_params(axis='x', rotation=45)
    
    # Add value labels on bars
    if any(r > 0 for r in avg_rps):
        min_rps = min(r for r in avg_rps if r > 0)
        for bar, value in zip(bars1, avg_rps):
            if value > 0:
                multiplier = value / min_rps
                ax1.text(bar.get_x() + bar.get_width()/2, bar.get_height() + max(avg_rps)*0.01,
                        f'{value:.0f} ({multiplier:.2f}x)', ha='center', va='bottom')
    
    # Average Mean Latency
    bars2 = ax2.bar(languages, avg_latency, color=colors)
    ax2.set_title('Average Mean Latency')
    ax2.set_ylabel('Latency (ms)')
    ax2.tick_params(axis='x', rotation=45)

    if any(l > 0 for l in avg_latency):
        min_latency = min(l for l in avg_latency if l > 0)
        for bar, value in zip(bars2, avg_latency):
            if value > 0:
                multiplier = value / min_latency
                ax2.text(bar.get_x() + bar.get_width()/2, bar.get_height() + max(avg_latency)*0.01,
                        f'{value:.2f}ms ({multiplier:.2f}x)', ha='center', va='bottom')

    # Average P95 Latency (changed from P99)
    bars3 = ax3.bar(languages, avg_p95_latency, color=colors)
    ax3.set_title('Average P95 Latency')  # Changed title from P99 to P95
    ax3.set_ylabel('P95 Latency (ms)')  # Changed ylabel from P99 to P95
    ax3.tick_params(axis='x', rotation=45)

    if any(l > 0 for l in avg_p95_latency):
        min_p95_latency = min(l for l in avg_p95_latency if l > 0)
        for bar, value in zip(bars3, avg_p95_latency):
            if value > 0:
                multiplier = value / min_p95_latency
                ax3.text(bar.get_x() + bar.get_width()/2, bar.get_height() + max(avg_p95_latency)*0.01,
                        f'{value:.2f}ms ({multiplier:.2f}x)', ha='center', va='bottom')

    # Average Memory
    bars4 = ax4.bar(languages, avg_memory, color=colors)
    ax4.set_title('Average Memory Usage')
    ax4.set_ylabel('Memory (MB)')
    ax4.tick_params(axis='x', rotation=45)

    if any(m > 0 for m in avg_memory):
        min_memory = min(m for m in avg_memory if m > 0)
        for bar, value in zip(bars4, avg_memory):
            if value > 0:
                multiplier = value / min_memory
                ax4.text(bar.get_x() + bar.get_width()/2, bar.get_height() + max(avg_memory)*0.01,
                        f'{value:.1f}MB ({multiplier:.2f}x)', ha='center', va='bottom')

    plt.tight_layout()
    plt.savefig(f'{output_dir}/summary_comparison.png', bbox_inches='tight')
    plt.close()

def main():
    parser = argparse.ArgumentParser(description='Compare benchmark results from multiple directories')
    parser.add_argument('directories', nargs='+', help='Directories containing stats.json and mem.json files')
    parser.add_argument('--output-dir', default='plots', help='Directory to save plots (default: plots)')
  
    args = parser.parse_args()
  
    dir_data_dict = {}
  
    for dir_path in args.directories:
        if not os.path.exists(dir_path):
            print(f"Warning: Directory {dir_path} not found, skipping...")
            continue
        
        if not os.path.isdir(dir_path):
            print(f"Warning: {dir_path} is not a directory, skipping...")
            continue
          
        try:
            stats_data, memory_data = load_directory_data(dir_path)
            dirname = Path(dir_path).name  # Get directory name
            dir_data_dict[dirname] = (stats_data, memory_data)
            
            stats_count = len(stats_data) if stats_data else 0
            memory_count = len(memory_data) if memory_data else 0
            print(f"Loaded {stats_count} stats points and {memory_count} memory points from {dir_path}")
        except Exception as e:
            print(f"Error loading {dir_path}: {e}")
  
    if not dir_data_dict:
        print("No valid data directories found!")
        return
  
    print(f"\nGenerating comparison plots for {len(dir_data_dict)} directories...")
    plot_comparisons(dir_data_dict, args.output_dir)
    print(f"Plots saved to {args.output_dir}/ directory")

if __name__ == "__main__":
    main()
