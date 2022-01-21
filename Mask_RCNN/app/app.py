from flask import Flask, Response, jsonify, request
import threading
from demo import run_demo
import base64
import numpy as np
import cv2
import re
import os
import json

app = Flask(__name__)

def convertImage(image):
    image = re.sub('data:.*?base64,', '', image)
    im_bytes = base64.b64decode(image)
    im_arr = np.frombuffer(im_bytes, dtype=np.uint8)
    return cv2.imdecode(im_arr, flags=cv2.IMREAD_COLOR)

@app.route('/upload_model', endpoint='upload_model', methods=['POST'])
def upload_model():
    model = request.files['model']
    MODELS_DIR = os.path.join(os.getcwd(),"models")
    model.save(os.path.join(MODELS_DIR, model.filename))
    return jsonify({'hello': 'MODELADDED', 'name': model.filename})

@app.route('/demo', endpoint='demo', methods=['POST'])
def demo():
    model = request.json.get('model')
    image = request.json.get('image')
    MODELS_DIR = os.path.join(os.getcwd(),"models")
    model = os.path.join(MODELS_DIR, model)
    id = request.json.get('id')
    run_demo(id, convertImage(image), model)
    return json.dumps({'id':id}), 200, {'ContentType':'application/json'}


if __name__ == '__main__':
    app.run(host="0.0.0.0", port=80)
    for threads in threading.enumerate():
        threads.join()
