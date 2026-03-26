(() => {
	const selector = '#QuoteForm input[name], #QuoteForm select[name], #QuoteForm textarea[name], #QuotePlans select[name], #QuoteForm button[name], #QuotePlans button[name]';
	const lastSent = new Map();
	const debounceMs = 250;
	const custNameDebounceMs = 450;
	const foldIds = ['QuoteInfoCard', 'QuoteSelectedCard'];
	const foldState = new Map();
	let seq = 0;
	let timer = 0;
	let pendingOpenSelected = false;
	let phoneStickyOn = false;
	let stickyFrame = 0;

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
			if (id === 'QuoteInfoCard') {
				el.open = true;
			}
		}
		if (pendingOpenSelected) {
			const selected = document.getElementById('QuoteSelectedCard');
			if (selected instanceof HTMLDetailsElement) {
				selected.open = true;
				foldState.set('QuoteSelectedCard', true);
			}
			pendingOpenSelected = false;
		}
	};

	const setPhoneSticky = (info, selected, on, collapse = true) => {
		phoneStickyOn = on;
		if (on) {
			info.classList.add('quote-phone-sticky-card');
			selected.classList.add('quote-phone-sticky-card', 'quote-phone-sticky-second');
			if (collapse) {
				if (info.open) {
					info.open = false;
					foldState.set('QuoteInfoCard', false);
				}
				if (selected.open) {
					selected.open = false;
					foldState.set('QuoteSelectedCard', false);
				}
			}
			return;
		}
		info.classList.remove('quote-phone-sticky-card');
		selected.classList.remove('quote-phone-sticky-card', 'quote-phone-sticky-second');
	};

	const syncPhoneSticky = () => {
		const info = document.getElementById('QuoteInfoCard');
		const selected = document.getElementById('QuoteSelectedCard');
		const anchor = document.getElementById('QuotePhoneStickyAnchor');
		if (!(info instanceof HTMLDetailsElement) || !(selected instanceof HTMLDetailsElement) || !(anchor instanceof HTMLElement)) {
			phoneStickyOn = false;
			return;
		}

		const top = anchor.getBoundingClientRect().top;
		if (!phoneStickyOn) {
			if (top <= 0) {
				setPhoneSticky(info, selected, true, true);
			}
			return;
		}

		// Hysteresis keeps sticky state stable during tiny back/forth scroll moves.
		if (top > 8) {
			setPhoneSticky(info, selected, false);
			return;
		}

		// After outerHTML rewrites, state may still be sticky but classes are gone on new nodes.
		if (!info.classList.contains('quote-phone-sticky-card') || !selected.classList.contains('quote-phone-sticky-card')) {
			setPhoneSticky(info, selected, true, false);
		}
	};

	const scheduleStickySync = () => {
		if (stickyFrame !== 0) return;
		stickyFrame = window.requestAnimationFrame(() => {
			stickyFrame = 0;
			syncPhoneSticky();
		});
	};

	const controlValue = (el) => {
		if (el instanceof HTMLButtonElement) return el.value || '1';
		return el.type === 'checkbox' ? (el.checked ? '1' : '0') : el.value;
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

	const postChange = (name, value) => {
		const call = ++seq;
		captureFoldStates();
		const form = new FormData();
		form.append('name', name);
		form.append('value', value);

		fetch('/quote-info-change', {
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
					if (msg.method === 'remove') {
						target.remove();
						continue;
					}
					if (msg.method === 'innerHTML') {
						target.innerHTML = msg.content;
						continue;
					}
					if (msg.method === 'outerHTML') {
						target.outerHTML = msg.content;
					}
				}
				applyFoldStates();
				syncPhoneSticky();
			})
			.catch(() => {});
	};

	const schedule = (name, value, wait) => {
		window.clearTimeout(timer);
		timer = window.setTimeout(() => {
			sendIfChanged(name, value);
		}, wait);
	};

	const onControlChange = (ev) => {
		const el = ev.target;
		if (!(el instanceof HTMLElement)) return;
		if (!el.matches(selector)) return;

		const name = el.getAttribute('name') || '';
		if (!name) return;
		const value = controlValue(el);
		if (name === 'custName') {
			schedule(name, value, custNameDebounceMs);
			return;
		}
		if (el.tagName === 'INPUT' && el.getAttribute('type') !== 'checkbox') {
			schedule(name, value, debounceMs);
			return;
		}
		sendIfChanged(name, value);
	};

	const onButtonClick = (ev) => {
		const el = ev.target;
		if (!(el instanceof HTMLElement)) return;
		if (!(el instanceof HTMLButtonElement)) return;
		if (!el.matches(selector)) return;
		const name = el.getAttribute('name') || '';
		if (!name) return;
		ev.preventDefault();
		if (name.startsWith('seladd-')) {
			pendingOpenSelected = true;
			foldState.set('QuoteSelectedCard', true);
		}
		const value = controlValue(el);
		sendIfChanged(name, value, true);
	};

	const onFoldToggle = (ev) => {
		const el = ev.target;
		if (!(el instanceof HTMLDetailsElement)) return;
		if (!foldIds.includes(el.id)) return;
		foldState.set(el.id, el.open);
	};

	document.addEventListener('change', onControlChange);
	document.addEventListener('input', onControlChange);
	document.addEventListener('click', onButtonClick);
	document.addEventListener('toggle', onFoldToggle, true);
	window.addEventListener('scroll', scheduleStickySync, { passive: true });
	window.addEventListener('resize', scheduleStickySync);
	scheduleStickySync();
})();
