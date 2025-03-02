import { toolTypes } from "../definitions"

//Utility function for drawing "rough.js" type elements on Canvas

export const drawBoardElement = ({roughCanvas, context, element}) => {

    switch (element.type) {
        case toolTypes.RECTANGLE:
            return roughCanvas.draw(element.roughElement);
        case toolTypes.CIRCLE:
            return roughCanvas.draw(element.roughElement);
        case toolTypes.LINE:
            return roughCanvas.draw(element.roughElement);
        default:
            throw new Error('An Error occured when drawing Element');
    }
} ;