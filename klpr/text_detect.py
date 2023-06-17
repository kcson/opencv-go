import sys

import cv2 as cv
import numpy as np

from imutils.object_detection import non_max_suppression

tm = cv.TickMeter()
tm.start()
src = cv.imread('../imgs/vehicle4.jpeg')
if src is None:
    print('image load fail!!')
    sys.exit()

origin = src.copy()
height, width, channels = src.shape
rW = width / float(320)
rH = height / float(320)

src = cv.resize(src, (320, 320))
height, width, channels = src.shape

output_layers = ['feature_fusion/Conv_7/Sigmoid',
                 'feature_fusion/concat_3']

net = cv.dnn.readNet('frozen_east_text_detection.pb')
blob = cv.dnn.blobFromImage(src, 1.0, (height, width), (123.68, 116.78, 103.94), swapRB=True, crop=False)
net.setInput(blob)

(scores, geometry) = net.forward(output_layers)
(numRows, numCols) = scores.shape[2:4]
rects = []
confidences = []

for y in range(0, numRows):
    scoresData = scores[0, 0, y]
    xData0 = geometry[0, 0, y]
    xData1 = geometry[0, 1, y]
    xData2 = geometry[0, 2, y]
    xData3 = geometry[0, 3, y]
    anglesData = geometry[0, 4, y]

    for x in range(0, numCols):
        # 만약 score가 충분한 확률을 가지고 있지 않다면 무시한다
        if scoresData[x] < 0.009:
            continue

        # 우리의 resulting feature map은 input_image보다 4배 작을것 이기 때문에
        # offset factor를 계산한다
        (offsetX, offsetY) = (x * 4.0, y * 4.0)

        # prediciton에 대한 회전각을 구하고 sin,cosine을 계산한다
        # 글씨가 회전되어 있을때를 대비
        angle = anglesData[x]
        cos = np.cos(angle)
        sin = np.sin(angle)

        # geometry volume를 사용해 bounding box의 width 와 height를 구한다
        h = xData0[x] + xData2[x]
        w = xData1[x] + xData3[x]

        # text prediction bounding box의 starting, ending (x,y) 좌표를 계산한다
        endX = int(offsetX + (cos * xData1[x]) + (sin * xData2[x]))
        endY = int(offsetY - (sin * xData1[x]) + (cos * xData2[x]))
        startX = int(endX - w)
        startY = int(endY - h)

        # bounding box coordinates와 probability score를 append한다
        rects.append((startX, startY, endX, endY))
        confidences.append(scoresData[x])

boxes = non_max_suppression(np.array(rects), probs=confidences)
for (startX, startY, endX, endY) in boxes:
    # 앞에서 구한 비율에 따라서 bounding box 좌표를 키워준다
    startX = int(startX * rW)
    startY = int(startY * rH)
    endX = int(endX * rW)
    endY = int(endY * rH)

    cv.rectangle(origin, (startX, startY), (endX, endY), (0, 255, 0), 2)

tm.stop()
print(tm.getTimeMilli())

cv.imshow('src', origin)
cv.waitKey()
cv.destroyAllWindows()
