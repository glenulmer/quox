(() => {
	const root = document.getElementById('QuoteInformation');
	if (!root) return;

	const selector = 'input[name], select[name], textarea[name]';
	const lastSent = new Map();
	const custNameDebounceMs = 450;
	let seq = 0;
	let custNameTimer = 0;

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

	const onControlChange = (ev) => {
		const el = ev.target;
		if (!(el instanceof HTMLElement)) return;
		if (!el.matches(selector)) return;

		const name = el.getAttribute('name') || '';
		if (!name) return;
		const value = controlValue(el);
		if (name === 'custName') {
			window.clearTimeout(custNameTimer);
			custNameTimer = window.setTimeout(() => {
				sendIfChanged(name, value);
			}, custNameDebounceMs);
			return;
		}
		sendIfChanged(name, value);
	};

	root.addEventListener('change', onControlChange);
	root.addEventListener('input', onControlChange);
})();
