export default class Utils {

    /** get object keys length */
    public static getObjectLength(o: object): number {
        if(!o || !(o instanceof Object)) {
            return 0
        }
        return Object.keys(o).length
    }
    
}