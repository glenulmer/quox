(() => {
	const onlyDigits = (v, max) => (v || '').replace(/\D+/g, '').slice(0, max);
	const pad2 = (n) => String(n).padStart(2, '0');

	const todayNum = () => {
		const d = new Date();
		return d.getFullYear() * 10000 + (d.getMonth() + 1) * 100 + d.getDate();
	};

	const daysInMonth = (year, month) => {
		if (month < 1 || month > 12) return 0;
		if (month === 2) {
			const leap = year % 400 === 0 || (year % 4 === 0 && year % 100 !== 0);
			return leap ? 29 : 28;
		}
		if (month === 4 || month === 6 || month === 9 || month === 11) return 30;
		return 31;
	};

	const parseISO = (v) => {
		const m = /^(\d{4})-(\d{2})-(\d{2})$/.exec((v || '').trim());
		if (!m) return null;
		const y = Number.parseInt(m[1], 10);
		const mm = Number.parseInt(m[2], 10);
		const dd = Number.parseInt(m[3], 10);
		if (dd < 1 || dd > daysInMonth(y, mm)) return null;
		return { y, mm, dd };
	};

	const isoNum = ({ y, mm, dd }) => y * 10000 + mm * 100 + dd;
	const isoText = ({ y, mm, dd }) => `${String(y).padStart(4, '0')}-${pad2(mm)}-${pad2(dd)}`;

	const validity = (wrap, msg) => {
		wrap.classList.toggle('qdate-invalid', !!msg);
		for (const part of wrap.querySelectorAll('input[data-qdate-part]')) {
			if (!(part instanceof HTMLInputElement)) continue;
			part.setCustomValidity(msg || '');
		}
	};

	const ruleFor = (wrap, hidden) => {
		const rule = (wrap.dataset.rule || '').trim();
		if (rule) return rule;
		const name = (hidden.name || '').toLowerCase();
		if (name === 'buy') return 'today-or-later';
		if (name === 'birth' || name.endsWith('-birth')) return 'past';
		return '';
	};

	const sync = (wrap, emit) => {
		const hidden = wrap.querySelector('input[data-qdate-hidden]');
		const day = wrap.querySelector('input[data-qdate-part="day"]');
		const month = wrap.querySelector('input[data-qdate-part="month"]');
		const year = wrap.querySelector('input[data-qdate-part="year"]');
		if (!(hidden instanceof HTMLInputElement) || !(day instanceof HTMLInputElement) || !(month instanceof HTMLInputElement) || !(year instanceof HTMLInputElement)) return;

		day.value = onlyDigits(day.value, 2);
		month.value = onlyDigits(month.value, 2);
		year.value = onlyDigits(year.value, 4);
		if (day.value.length !== 2 || month.value.length !== 2 || year.value.length !== 4) {
			validity(wrap, '');
			return;
		}

		const parsed = parseISO(`${year.value}-${month.value}-${day.value}`);
		if (!parsed) {
			validity(wrap, 'Invalid date.');
			return;
		}

		const n = isoNum(parsed);
		const rule = ruleFor(wrap, hidden);
		if (rule === 'past' && n >= todayNum()) {
			validity(wrap, 'Date must be in the past.');
			return;
		}
		if (rule === 'today-or-later' && n < todayNum()) {
			validity(wrap, 'Date must be today or later.');
			return;
		}
		validity(wrap, '');

		const iso = isoText(parsed);
		if (hidden.value === iso) return;
		hidden.value = iso;
		if (emit) {
			hidden.dispatchEvent(new Event('input', { bubbles: true }));
			hidden.dispatchEvent(new Event('change', { bubbles: true }));
		}
	};

	const hydrate = (wrap) => {
		const hidden = wrap.querySelector('input[data-qdate-hidden]');
		const day = wrap.querySelector('input[data-qdate-part="day"]');
		const month = wrap.querySelector('input[data-qdate-part="month"]');
		const year = wrap.querySelector('input[data-qdate-part="year"]');
		if (!(hidden instanceof HTMLInputElement) || !(day instanceof HTMLInputElement) || !(month instanceof HTMLInputElement) || !(year instanceof HTMLInputElement)) return;
		const parsed = parseISO(hidden.value);
		if (!parsed) return;
		day.value = pad2(parsed.dd);
		month.value = pad2(parsed.mm);
		year.value = String(parsed.y).padStart(4, '0');
	};

	const wire = (wrap) => {
		if (!(wrap instanceof HTMLElement) || wrap.dataset.qdateReady === '1') return;
		wrap.dataset.qdateReady = '1';
		hydrate(wrap);
		for (const part of wrap.querySelectorAll('input[data-qdate-part]')) {
			if (!(part instanceof HTMLInputElement)) continue;
			part.addEventListener('input', () => {
				const role = part.dataset.qdatePart || '';
				const max = role === 'year' ? 4 : 2;
				part.value = onlyDigits(part.value, max);
				if ((role === 'day' || role === 'month') && part.value.length === 2) {
					const next = role === 'day'
						? wrap.querySelector('input[data-qdate-part="month"]')
						: wrap.querySelector('input[data-qdate-part="year"]');
					if (next instanceof HTMLInputElement) next.focus();
				}
				sync(wrap, true);
			});
			part.addEventListener('blur', () => sync(wrap, true));
		}
	};

	window.QuoDateParts = {
		init(root = document) {
			for (const wrap of root.querySelectorAll('.qdate[data-qdate="1"]')) wire(wrap);
		},
	};
})();
