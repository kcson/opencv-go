import sys

import cv2 as cv
import numpy as np

tm = cv.TickMeter()
tm.start()
net = cv.dnn.readNetFromDarknet('yolov4-ANPR.cfg', 'yolov4-ANPR.weights')
# net = cv.dnn.readNet('yolov4-ANPR.weights', 'yolov4-ANPR.cfg')
# classes = []
with open('obj.names', 'r') as f:
    classes = [line.strip() for line in f.readlines()]

layer_names = net.getLayerNames()
output_layers = [layer_names[i - 1] for i in net.getUnconnectedOutLayers()]
colors = np.random.uniform(0, 255, size=(len(classes), 3))

src = cv.imread('../imgs/vehicle1.jpeg')
if src is None:
    print('image load fail!!')
    sys.exit()

height, width, channels = src.shape
origin = src.copy()
src = cv.resize(src,  None, fx=416/width, fy=416/height)
blob = cv.dnn.blobFromImage(src, 0.00392, (416, 416), (0, 0, 0), True, crop=False)
net.setInput(blob)

outs = net.forward(output_layers)
class_ids = []
confidences = []
boxes = []
for out in outs:
    for detection in out:
        scores = detection[5:]
        class_id = np.argmax(scores)
        confidence = scores[class_id]

        if confidence > 0.75:
            center_x = int(detection[0] * width)
            center_y = int(detection[1] * height)
            w = int(detection[2] * width)
            h = int(detection[3] * height)

            # 객체의 사각형 테두리 중 좌상단 좌표값 찾기
            x = int(center_x - w / 2)
            y = int(center_y - h / 2)

            boxes.append([x, y, w, h])
            confidences.append(float(confidence))
            class_ids.append(class_id)

indexes = cv.dnn.NMSBoxes(boxes, confidences, 0.5, 0.4)
max_box_area = 0
max_box = (0, 0, 0, 0)
for i in range(len(boxes)):
    if i in indexes:
        class_name = classes[class_ids[i]]
        if class_name == 'car' and len(indexes) > 1:
            continue
        x, y, w, h = boxes[i]
        if w * h > max_box_area:
            max_box = boxes[i]
        # cv.rectangle(src, (x, y, w, h), (0, 0, 255), thickness=3)

x, y, w, h = max_box
cv.rectangle(origin, (x, y, w, h), (0, 0, 255), thickness=3)
tm.stop()
print(tm.getTimeMilli())

cv.imshow('src', origin)
cv.waitKey()
cv.destroyAllWindows()
