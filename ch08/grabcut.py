import sys
import numpy as np
import cv2 as cv

# 입력 영상 불러오기
src = cv.imread('../imgs/nemo.jpg')

if src is None:
    print('Image load failed!')
    sys.exit()

# 사각형 지정을 통한 초기 분할
rc = cv.selectROI(src)
mask = np.zeros(src.shape[:2], np.uint8)

cv.grabCut(src, mask, rc, None, None, 5, cv.GC_INIT_WITH_RECT)

# 0: cv2.GC_BGD, 2: cv2.GC_PR_BGD
mask2 = np.where((mask == 0) | (mask == 2), 0, 1).astype('uint8')

cv.imshow('mask2', mask2 * 255)
dst = src * mask2[:, :, np.newaxis]

# 초기 분할 결과 출력
cv.imshow('mask', mask * 64)
cv.imshow('dst', dst)
cv.waitKey()
cv.destroyAllWindows()