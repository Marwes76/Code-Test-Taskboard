export class OrderBy {
	static readonly SORT_ORDER =	new OrderBy("SORT_ORDER");
	static readonly ALPHABETICAL =	new OrderBy("ALPHABETICAL");
	static readonly CREATED_AT =	new OrderBy("CREATED_AT");
	static readonly UPDATED_AT =	new OrderBy("UPDATED_AT");

	private constructor(public readonly value: string) {}

	isSortOrder(): boolean {
		return this === OrderBy.SORT_ORDER;
	}

	isAlphabetical(): boolean {
		return this === OrderBy.ALPHABETICAL;
	}

	isCreatedAt(): boolean {
		return this === OrderBy.CREATED_AT;
	}

	isUpdatedAt(): boolean {
		return this === OrderBy.UPDATED_AT;
	}

	static toggleOrderBy(orderBy: OrderBy): OrderBy {
		if (orderBy.isSortOrder()) {
			return OrderBy.ALPHABETICAL;
		} else if (orderBy.isAlphabetical()) {
			return OrderBy.CREATED_AT;
		} else if (orderBy.isCreatedAt()) {
			return OrderBy.UPDATED_AT;
		} else if (orderBy.isUpdatedAt()) {
			return OrderBy.SORT_ORDER;
		}

		return OrderBy.ALPHABETICAL;
	}

	getOrderByIcon() {
		if (this.isSortOrder()) {
			return "format_list_numbered";
		} else if (this.isAlphabetical()) {
			return "sort_by_alpha";
		} else if (this.isCreatedAt()) {
			return "access_time";
		} else if (this.isUpdatedAt()) {
			return "edit";
		}

		return "";
	}

	getOrderByString() {
		if (this.isSortOrder()) {
			return "Sort order";
		} else if (this.isAlphabetical()) {
			return "Alphabetical";
		} else if (this.isCreatedAt()) {
			return "Created at";
		} else if (this.isUpdatedAt()) {
			return "Updated at";
		}

		return "";
	}
}