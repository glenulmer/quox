(() => {
	const selector = '#QuoteForm input[name], #QuoteForm select[name], #QuoteForm textarea[name], #QuotePlans select[name^="plancat-"]';
	const lastSent = new Map();
	const debounceMs = 250;
	const custNameDebounceMs = 450;
	let seq = 0;
	let timer = 0;

	const controlValue = (el) => (el.type === 'checkbox' ? (el.checked ? '1' : '0') : el.value);
	const sendIfChanged = (name, value) => {
		if (lastSent.get(name) === value) return;
		lastSent.set(name, value);
		postChange(name, value);
	};

	const postChange = (name, value) => {
		const call = ++seq;
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

	document.addEventListener('change', onControlChange);
	document.addEventListener('input', onControlChange);
})();
