import rough from 'roughjs/bundled/rough.esm';
import { toolTypes } from '../definitions';

const generator = rough.generator();

const generateRectangle = ({x1, y1, x2, y2}) => {
    return generator.rectangle(x1, y1, x2 - x1, y2 - y1);

};

const generateCircle = ({x1, y1, x2, y2}) => {
    //console.log("x1:", x1, "y1:", y1); console.log("x2:", x2, "y2:", y2); 
    //var radius = 10 * Math.sqrt((x2 - x1)^2 + (y2 - y1)^2);
    var radius = ((x2 - x1) + (y2 - y1));
    //console.log("Radius:", radius);
    return generator.circle(x1, y1, radius);
  };

const generateLine = ({x1, y1, x2, y2}) => {
    return generator.line(x1, y1, x2, y2);

};

export const  createBoardElement = ({ x1, y1, x2, y2, toolType, id}) => {
    let roughElement;

    switch (toolType) {
        case toolTypes.RECTANGLE:
            roughElement = generateRectangle({x1, y1, x2, y2});
            return {
                id: id,
                roughElement,
                type:toolType,
                x1,
                y1,
                x2,
                y2,

            };
        case toolTypes.CIRCLE:
            roughElement = generateCircle({x1, y1, x2, y2});
            return {
                id: id,
                roughElement,
                type:toolType,
                x1,
                y1,
                x2,
                y2,
            };
        case toolTypes.LINE:
            roughElement = generateLine({x1, y1, x2, y2});
            return {
                    id: id,
                    roughElement,
                    type:toolType,
                    x1,
                    y1,
                    x2,
                    y2,
                };  
           
        default:
                throw new Error("Unexpected Error on creating new element");
    }
};