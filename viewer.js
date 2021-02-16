const imageSize = 1000;
const server = "http://localhost:9000"

class Viewport {
    constructor(topLeft, bottomRight) {
        this.topLeft = topLeft;
        this.bottomRight = bottomRight;
    }

    width() {
        return this.bottomRight.x - this.topLeft.x
    }

    height() {
        return this.topLeft.y - this.bottomRight.y
    }

    pointAt(event) {
        const xDelta = event.offsetX / imageSize * this.width();
        const yDelta = event.offsetY / imageSize * this.height();
        const realPart = this.topLeft.x + xDelta;
        const imagPart = this.topLeft.y - yDelta;
        console.log("point is", realPart, imagPart)
        return {x: realPart, y: imagPart}
    }

    zoomIn(center) {
        this.updateCenter(center, 2/3);
    }

    zoomOut(center) {
        this.updateCenter(center, 3/2);
    }

    updateCenter(center, zoomFactor) {
        const w = (this.width() / 2) * zoomFactor;
        this.topLeft = {
            x: center.x - w,
            y: center.y + w
        }
        this.bottomRight = {
            x: center.x + w,
            y: center.y - w
        }
        console.log("viewport width is", this.width());
    }
}

$(() => {
    const viewport = new Viewport({x: -2, y: 2}, {x: 2, y: -2});
    const viewportElem = $("#viewport");

    viewportElem.on("wheel", event => {
        console.log("got event", {
            offsetX: event.offsetX,
            offsetY: event.offsetY,
            wheelDelta: event.originalEvent.wheelDelta
        });
        const center = viewport.pointAt(event);
        if (event.originalEvent.wheelDelta > 0) {
            viewport.zoomIn(center)
        } else {
            viewport.zoomOut(center);
        }

        const c = encodeURIComponent(`${center.x}+${center.y}i`);
        const w = viewport.width();
        const url = `${server}?center=${c}&width=${w}`;
        viewportElem.attr("src", url);
        return false;
    });
})
