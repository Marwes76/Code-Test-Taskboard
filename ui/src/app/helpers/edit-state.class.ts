export class EditState {
	static readonly NEW =		new EditState("NEW");
	static readonly DEFAULT =	new EditState("DEFAULT");
	static readonly EDITING =	new EditState("EDITING");
	static readonly SAVING =	new EditState("SAVING");
	static readonly LOCKED =	new EditState("LOCKED"); // Locked for other reasons; for example, you might want to lock parent when editing a child

	private constructor(public readonly value: string) {}

	isNew(): boolean {
		return this === EditState.NEW;
	}

	isDefault(): boolean {
		return this === EditState.DEFAULT;
	}

	isEditing(): boolean {
		return this === EditState.EDITING;
	}

	isSaving(): boolean {
		return this === EditState.SAVING;
	}

	isLocked(): boolean {
		return this === EditState.LOCKED;
	}
}