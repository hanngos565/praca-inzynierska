import os
import threading
import coco
import model as modellib
from keras import backend as K
from update import update

def fire_and_forget(f):
    def wrapped(*args, **kwargs):
        threading.Thread(target=f, args=args, kwargs=kwargs).start()
    return wrapped

@fire_and_forget
def run_demo(id, image, model):
    try:
        ROOT_DIR = os.getcwd()
        MODEL_DIR = os.path.join(ROOT_DIR, "logs")
        COCO_MODEL_PATH = model

        K.clear_session()

        class InferenceConfig(coco.CocoConfig):
            GPU_COUNT = 1
            IMAGES_PER_GPU = 1

        config = InferenceConfig()

        model = modellib.MaskRCNN(mode="inference", model_dir=MODEL_DIR, config=config)

        model.load_weights(COCO_MODEL_PATH, by_name=True)

        class_names = ['BG', 'person', 'bicycle', 'car', 'motorcycle', 'airplane',
                       'bus', 'train', 'truck', 'boat', 'traffic light',
                       'fire hydrant', 'stop sign', 'parking meter', 'bench', 'bird',
                       'cat', 'dog', 'horse', 'sheep', 'cow', 'elephant', 'bear',
                       'zebra', 'giraffe', 'backpack', 'umbrella', 'handbag', 'tie',
                       'suitcase', 'frisbee', 'skis', 'snowboard', 'sports ball',
                       'kite', 'baseball bat', 'baseball glove', 'skateboard',
                       'surfboard', 'tennis racket', 'bottle', 'wine glass', 'cup',
                       'fork', 'knife', 'spoon', 'bowl', 'banana', 'apple',
                       'sandwich', 'orange', 'broccoli', 'carrot', 'hot dog', 'pizza',
                       'donut', 'cake', 'chair', 'couch', 'potted plant', 'bed',
                       'dining table', 'toilet', 'tv', 'laptop', 'mouse', 'remote',
                       'keyboard', 'cell phone', 'microwave', 'oven', 'toaster',
                       'sink', 'refrigerator', 'book', 'clock', 'vase', 'scissors',
                       'teddy bear', 'hair drier', 'toothbrush']

        results = model.detect([image], verbose=1)

        r = results[0]

        def ids_to_names(n):
            return class_names[n]

        class_ids_to_names = map(ids_to_names, r['class_ids'])

        class_ids_to_names = list(class_ids_to_names)


        results = {
            'names': class_ids_to_names,
            'scores': r['scores'].tolist(),
            'bbox': r['rois'].tolist(),
        }

    except:
        results = "error"
    finally:
        K.clear_session()
        update(id, results)
