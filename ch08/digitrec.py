import sys
import numpy as np
import cv2 as cv


def load_digits():
    img_digits = []

    for i in range(10):
        filename = '../imgs/digits/digit{}.bmp'.format(i)
        img_digits.append(cv.imread(filename, cv.IMREAD_GRAYSCALE))

        if img_digits[i] is None:
            return None
    return img_digits


def find_digit(img, img_digits):
    max_idx = -1
    max_ccoeff = -1

    img = cv.resize(img, (100, 150))
    for i in range(10):
        res = cv.matchTemplate(img, img_digits[i], cv.TM_CCOEFF_NORMED)
        if res[0, 0] > max_ccoeff:
            max_idx = i
            max_ccoeff = res[0, 0]
    return max_idx


def main():
    src = cv.imread("../imgs/digits_print.bmp")
    # src = cv.imread("../imgs/vehicle.jpeg")
    if src is None:
        print("image load fail!!")
        sys.exit()

    img_digits = load_digits()
    if img_digits is None:
        print('Digits load fail!')
        return

    src_gray = cv.cvtColor(src, cv.COLOR_BGR2GRAY)
    ret, src_bin = cv.threshold(src_gray, 0, 255, cv.THRESH_BINARY_INV | cv.THRESH_OTSU)
    cnt, _, stats, _ = cv.connectedComponentsWithStats(src_bin)
    cv.imshow('src_bin', src_bin)

    dst = src.copy()
    for i in range(1, cnt):
        (x, y, w, h, s) = stats[i]
        if s < 1000:
            continue
        digit = find_digit(src_gray[y:y + h, x:x + w], img_digits)
        cv.rectangle(dst, (x, y, w, h), (0, 255, 255))
        cv.putText(dst, str(digit), (x, y - 4), cv.FONT_HERSHEY_SIMPLEX, 1, (0, 255, 255), 2, cv.LINE_AA)

    cv.imshow('dst', dst)

    cv.waitKey()
    cv.destroyAllWindows()


if __name__ == "__main__":
    main()
