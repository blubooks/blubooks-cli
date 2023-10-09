export default BreakToken;
/**
 * BreakToken
 * @class
 */
declare class BreakToken {
    constructor(node: any, offset: any);
    node: any;
    offset: any;
    equals(otherBreakToken: any): boolean;
    toJSON(hash: any): {};
}
