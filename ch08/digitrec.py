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


def main():
    src = cv.imread("../imgs/digits_print.bmp")
    if src is None:
        print("image load fail!!")
        sys.exit()

    cv.imshow('src', src)

    cv.waitKey()
    cv.destroyAllWindows()


if __name__ == "__main__":
    main()
