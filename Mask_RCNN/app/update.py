import requests
import json

address = 'http://server:8081'

def update(id, results):
    results = json.dumps(results)
    data = json.dumps({'id': id, 'content': results})
    requests.put(address + '/v1/simulation-results', data=data, headers={'Content-Type': 'application/json'})