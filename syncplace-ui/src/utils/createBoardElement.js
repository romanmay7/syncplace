import rough from 'roughjs/bundled/rough.esm';
import { toolTypes } from '../definitions';

//Rough factory instance for generating different shapes 
const generator = rough.generator();

//Helper Functions for crating  elements of specific shapes
//***************************************************************************************** */
const generateRectangle = ({x1, y1, x2, y2, colour, fillMode}) => {
    if(!fillMode)
      return generator.rectangle(x1, y1, x2 - x1, y2 - y1, {stroke: colour});
    else
      return generator.rectangle(x1, y1, x2 - x1, y2 - y1, {fill: colour});

};
//-----------------------------------------------------------------------------------------

const generateCircle = ({x1, y1, x2, y2, colour, fillMode}) => {
    //console.log("x1:", x1, "y1:", y1); console.log("x2:", x2, "y2:", y2); 
    //var radius = 10 * Math.sqrt((x2 - x1)^2 + (y2 - y1)^2);
    var radius = ((x2 - x1) + (y2 - y1));
    //console.log("Radius:", radius);
    if(!fillMode)
      return generator.circle(x1, y1, radius, {stroke: colour});
    else
      return generator.circle(x1, y1, radius, {fill: colour});
  };
//-----------------------------------------------------------------------------------------

const generateLine = ({x1, y1, x2, y2, colour}) => {
    return generator.line(x1, y1, x2, y2, {stroke: colour});

};
//***************************************************************************************** */

// Common Utility function for creating elements of different shapes
export const  createBoardElement = ({ x1, y1, x2, y2, toolType, colour, fillMode, id}) => {
    let roughElement;

    switch (toolType) {
        case toolTypes.RECTANGLE:
            roughElement = generateRectangle({x1, y1, x2, y2, colour, fillMode});
            return {
                id: id,
                roughElement,
                type:toolType,
                colour:colour,
                x1,
                y1,
                x2,
                y2,

            };
        case toolTypes.CIRCLE:
            roughElement = generateCircle({x1, y1, x2, y2, colour, fillMode});
            return {
                id: id,
                roughElement,
                type:toolType,
                colour:colour,
                x1,
                y1,
                x2,
                y2,
            };
        case toolTypes.LINE:
            roughElement = generateLine({x1, y1, x2, y2, colour});
            return {
                    id: id,
                    roughElement,
                    type:toolType,
                    colour:colour,
                    x1,
                    y1,
                    x2,
                    y2,
                };  
           
        default:
                throw new Error("Unexpected Error on creating new element");
    }
};
//***************************************************************************************** */