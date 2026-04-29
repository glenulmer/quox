(() => {
	const selector = '#EditQForm input[name], #EditQForm textarea[name], #EditQForm select[name], #EditQForm button[name]';
	const lastSent = new Map();
	const debounceMs = 240;
	const foldIds = ['EditQPreexCard', 'EditQDependantsCard', 'EditQReviewCard'];
	const foldDefaults = new Map([
		['EditQPreexCard', false],
		['EditQDependantsCard', false],
		['EditQReviewCard', true],
	]);
	const foldState = new Map();
	let seq = 0;
	let timer = 0;

	const captureFoldStates = () => {
		for (const id of foldIds) {
			const el = document.getElementById(id);
			if (!(el instanceof HTMLDetailsElement)) continue;
			foldState.set(id, el.open);
		}
	};

	const applyFoldStates = () => {
		for (const id of foldIds) {
			const el = document.getElementById(id);
			if (!(el instanceof HTMLDetailsElement)) continue;
			if (foldState.has(id)) {
				el.open = !!foldState.get(id);
				continue;
			}
			if (foldDefaults.has(id)) {
				el.open = !!foldDefaults.get(id);
			}
		}
	};

	const captureFocus = () => {
		const active = document.activeElement;
		if (!(active instanceof HTMLElement)) return null;
		if (!(active instanceof HTMLInputElement) && !(active instanceof HTMLTextAreaElement)) return null;
		const name = active.getAttribute('name') || '';
		if (!name) return null;
		const out = { name, start: 0, end: 0 };
		if (typeof active.selectionStart === 'number') out.start = active.selectionStart;
		if (typeof active.selectionEnd === 'number') out.end = active.selectionEnd;
		return out;
	};

	const restoreFocus = (focus) => {
		if (!focus) return;
		const list = document.querySelectorAll('#EditQForm [name]');
		let target = null;
		for (const el of list) {
			if (!(el instanceof HTMLInputElement) && !(el instanceof HTMLTextAreaElement)) continue;
			if ((el.getAttribute('name') || '') !== focus.name) continue;
			target = el;
			break;
		}
		if (!target) return;
		target.focus();
		if (typeof target.setSelectionRange === 'function') {
			target.setSelectionRange(focus.start, focus.end);
		}
	};

	const autosize = (el) => {
		if (!(el instanceof HTMLTextAreaElement)) return;
		if (!el.classList.contains('editq-grow-input')) return;
		el.style.height = 'auto';
		el.style.height = `${Math.max(el.scrollHeight, 34)}px`;
	};

	const autosizeAll = () => {
		const list = document.querySelectorAll('#EditQForm textarea.editq-grow-input');
		for (const el of list) autosize(el);
	};

	const controlValue = (el) => {
		if (el instanceof HTMLButtonElement) return el.value || '1';
		if (el instanceof HTMLInputElement && el.type === 'checkbox') return el.checked ? '1' : '0';
		return el.value;
	};

	const postChange = (name, value) => {
		const call = ++seq;
		const focus = captureFocus();
		captureFoldStates();
		const form = new FormData();
		form.append('name', name);
		form.append('value', value);

		fetch('/quote-review-change', {
			method: 'POST',
			body: form,
			credentials: 'same-origin',
		})
			.then((res) => (res.ok ? res.json() : []))
			.then((messages) => {
				if (call !== seq || !Array.isArray(messages)) return;
				for (const msg of messages) {
					if (!msg || msg.kind !== 'rewrite') continue;
					const target = document.querySelector(msg.target);
					if (!target) continue;
					if (msg.method === 'outerHTML') target.outerHTML = msg.content;
					if (msg.method === 'innerHTML') target.innerHTML = msg.content;
					if (msg.method === 'remove') target.remove();
					}
					applyFoldStates();
					autosizeAll();
					restoreFocus(focus);
				})
			.catch(() => {});
	};

	const sendIfChanged = (name, value, force = false) => {
		if (force) {
			postChange(name, value);
			return;
		}
		if (lastSent.get(name) === value) return;
		lastSent.set(name, value);
		postChange(name, value);
	};

	const schedule = (name, value) => {
		window.clearTimeout(timer);
		timer = window.setTimeout(() => sendIfChanged(name, value), debounceMs);
	};

	const onChange = (ev) => {
		const el = ev.target;
		if (!(el instanceof HTMLElement)) return;
		if (!el.matches(selector)) return;
		if (el instanceof HTMLButtonElement) return;
		const name = el.getAttribute('name') || '';
		if (!name) return;
		if (name === 'lang' || name === 'slim') return;
		const value = controlValue(el);
		if (el instanceof HTMLTextAreaElement) {
			autosize(el);
			schedule(name, value);
			return;
		}
		if (el instanceof HTMLInputElement && el.type === 'text') {
			schedule(name, value);
			return;
		}
		sendIfChanged(name, value);
	};

	const onClick = (ev) => {
		const el = ev.target;
		if (!(el instanceof HTMLButtonElement)) return;
		if (!el.matches(selector)) return;
		const name = el.getAttribute('name') || '';
		if (!name) return;
		if (name === 'DownloadExcel') return;
		ev.preventDefault();
		sendIfChanged(name, controlValue(el), true);
	};

	const onFoldToggle = (ev) => {
		const el = ev.target;
		if (!(el instanceof HTMLDetailsElement)) return;
		if (!foldIds.includes(el.id)) return;
		foldState.set(el.id, el.open);
	};

	document.addEventListener('input', onChange);
	document.addEventListener('change', onChange);
	document.addEventListener('click', onClick);
	document.addEventListener('toggle', onFoldToggle, true);
	applyFoldStates();
	autosizeAll();
})();
