import { toolTypes } from "../definitions"


export const drawBoardElement = ({roughCanvas, context, element}) => {

    switch (element.type) {
        case toolTypes.RECTANGLE:
            return roughCanvas.draw(element.roughElement);
        default:
            throw new Error('An Error occured when drawing Element');
    }
} ;