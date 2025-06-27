from flask import Flask, jsonify
import json

app = Flask(__name__)

@app.route('/user/<int:id>')
def get_user(id):
    with open('main.json') as f:
        users = json.load(f)

    user = next((user for user in users if user['id'] == id), None)
    if user:
        return jsonify(user)
    return jsonify({'error': 'User not found'}), 404

if __name__ == '__main__':
    app.run(port=3003)
