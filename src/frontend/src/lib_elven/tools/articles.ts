export default class ToolsArticles {
	/** validate id */
	public static validateID(val: string): boolean {
		return !!val;
	}

	/** validate title */
	public static validateTitle(val: string): boolean {
		return val.length <= 124;
	}

	/** validate content */
	public static validateContent(val: string): boolean {
		return val && val.length <= 256000;
	}
}
