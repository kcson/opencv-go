from __future__ import annotations
from abc import ABC, abstractmethod
import cv2 as cv
import sys


class Handler(ABC):
    @abstractmethod
    def set_next(self, handler: Handler) -> Handler:
        pass

    @abstractmethod
    def handle(self, param) -> dict:
        pass


class AbstractHandler(Handler):
    _next_handler: Handler = None

    def set_next(self, handler: Handler) -> Handler:
        self._next_handler = handler

        return handler

    def handle(self, param) -> dict:
        if self._next_handler:
            return self._next_handler.handle(param)

        return None


if __name__ == "__main__":
    print("main")
    src = cv.imread('../imgs/vehicle.jpeg')
    if src is None:
        print('image read fail!!')
        sys.exit()
    cv.imshow('src', src)

    cv.waitKey()
    cv.destroyAllWindows()
