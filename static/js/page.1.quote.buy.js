(() => {
	const MONTHS = [`Jan`, `Feb`, `Mar`, `Apr`, `May`, `Jun`, `Jul`, `Aug`, `Sep`, `Oct`, `Nov`, `Dec`];
	const DAYS = [`Mo`, `Tu`, `We`, `Th`, `Fr`, `Sa`, `Su`];
	let popup = null;
	let monthSel = null;
	let yearInp = null;
	let grid = null;
	let wrap = null;
	let hidden = null;
	let trigger = null;
	let viewYear = 0;
	let viewMonth = 0;
	let selected = 0;
	let min = 0;
	let max = 0;

	const pad2 = (n) => String(n).padStart(2, `0`);
	const ymd = (y, m, d) => y * 10000 + m * 100 + d;
	const daysInMonth = (y, m) => new Date(Date.UTC(y, m, 0)).getUTCDate();
	const leadDays = (y, m) => (new Date(Date.UTC(y, m - 1, 1)).getUTCDay() + 6) % 7;
	const today = () => {
		const d = new Date();
		return ymd(d.getFullYear(), d.getMonth() + 1, d.getDate());
	};
	const inRange = (n) => (!min || n >= min) && (!max || n <= max);
	const monthAllowed = (y, m) => {
		const first = ymd(y, m, 1);
		const last = ymd(y, m, daysInMonth(y, m));
		if (min && last < min) return false;
		if (max && first > max) return false;
		return true;
	};
	const shiftedMonth = (y, m, step) => {
		let nm = m + step;
		let ny = y;
		if (nm < 1) { nm = 12; ny--; }
		if (nm > 12) { nm = 1; ny++; }
		return { y: ny, m: nm };
	};
	const shiftMonth = (step) => {
		const next = shiftedMonth(viewYear, viewMonth, step);
		if (!monthAllowed(next.y, next.m)) return;
		viewYear = next.y;
		viewMonth = next.m;
	};

	const parse8 = (v) => {
		v = (v || ``).trim();
		if (!/^\d{8}$/.test(v)) return 0;
		const y = Number.parseInt(v.slice(0, 4), 10);
		const m = Number.parseInt(v.slice(4, 6), 10);
		const d = Number.parseInt(v.slice(6, 8), 10);
		if (m < 1 || m > 12) return 0;
		if (d < 1 || d > daysInMonth(y, m)) return 0;
		return ymd(y, m, d);
	};

	const as8 = (n) => {
		const y = Math.trunc(n / 10000);
		const m = Math.trunc((n % 10000) / 100);
		const d = n % 100;
		return `${String(y).padStart(4, `0`)}${pad2(m)}${pad2(d)}`;
	};

	const asDMY = (n) => {
		const y = Math.trunc(n / 10000);
		const m = Math.trunc((n % 10000) / 100);
		const d = n % 100;
		return `${pad2(d)}.${pad2(m)}.${String(y).padStart(4, `0`)}`;
	};

	const setPopupPos = () => {
		if (!popup || !wrap) return;
		if (!document.body.contains(wrap)) { closePopup(); return; }
		const r = wrap.getBoundingClientRect();
		popup.style.left = `${Math.round(r.left)}px`;
		popup.style.top = `${Math.round(r.bottom + 4)}px`;
		popup.style.minWidth = `${Math.max(180, Math.round(r.width))}px`;
	};

	const closePopup = () => {
		if (!popup) return;
		popup.classList.remove(`is-open`);
		wrap = null;
		hidden = null;
		trigger = null;
	};

	const render = () => {
		if (!popup || !monthSel || !yearInp || !grid) return;
		monthSel.value = String(viewMonth);
		yearInp.value = String(viewYear);
		const prevBtn = popup.querySelector(`.qbuy-nav[data-step="-1"]`);
		const nextBtn = popup.querySelector(`.qbuy-nav[data-step="1"]`);
		const prev = shiftedMonth(viewYear, viewMonth, -1);
		const next = shiftedMonth(viewYear, viewMonth, 1);
		if (prevBtn instanceof HTMLButtonElement) prevBtn.disabled = !monthAllowed(prev.y, prev.m);
		if (nextBtn instanceof HTMLButtonElement) nextBtn.disabled = !monthAllowed(next.y, next.m);
		grid.innerHTML = ``;

		for (const d of DAYS) {
			const cell = document.createElement(`div`);
			cell.className = `qbuy-dow`;
			cell.textContent = d;
			grid.appendChild(cell);
		}

		for (let i = 0; i < leadDays(viewYear, viewMonth); i++) {
			const blank = document.createElement(`div`);
			blank.className = `qbuy-blank`;
			grid.appendChild(blank);
		}

		const t = today();
		const maxDay = daysInMonth(viewYear, viewMonth);
		for (let d = 1; d <= maxDay; d++) {
			const n = ymd(viewYear, viewMonth, d);
			const btn = document.createElement(`button`);
			btn.type = `button`;
			btn.className = `qbuy-day`;
			btn.dataset.day = String(d);
			btn.textContent = String(d);
			if (n === selected) btn.classList.add(`is-selected`);
			if (n === t) btn.classList.add(`is-today`);
			if (!inRange(n)) {
				btn.disabled = true;
				btn.classList.add(`is-disabled`);
			}
			grid.appendChild(btn);
		}
	};

	const ensurePopup = () => {
		if (popup) return;
		popup = document.createElement(`div`);
		popup.className = `qbuy-popup`;
		popup.innerHTML = `
			<div class="qbuy-popup-head">
				<button type="button" class="qbuy-nav" data-step="-1">‹</button>
				<div class="qbuy-head-mid">
					<select class="qbuy-month"></select>
					<input class="qbuy-year" type="text" inputmode="numeric" maxlength="4">
				</div>
				<button type="button" class="qbuy-nav" data-step="1">›</button>
			</div>
			<div class="qbuy-grid"></div>
		`;
		document.body.appendChild(popup);
		monthSel = popup.querySelector(`.qbuy-month`);
		yearInp = popup.querySelector(`.qbuy-year`);
		grid = popup.querySelector(`.qbuy-grid`);

		for (let i = 0; i < MONTHS.length; i++) {
			const opt = document.createElement(`option`);
			opt.value = String(i + 1);
			opt.textContent = MONTHS[i];
			monthSel.appendChild(opt);
		}

		monthSel.addEventListener(`change`, () => {
			const month = Number.parseInt(monthSel.value || `0`, 10);
			if (month < 1 || month > 12) return;
			viewMonth = month;
			render();
		});

		yearInp.addEventListener(`input`, () => {
			yearInp.value = (yearInp.value || ``).replace(/\D+/g, ``).slice(0, 4);
		});
		yearInp.addEventListener(`change`, () => {
			const year = Number.parseInt(yearInp.value || `0`, 10);
			if (!Number.isFinite(year) || year <= 0) {
				yearInp.value = String(viewYear);
				return;
			}
			const minYear = min ? Math.trunc(min / 10000) : 0;
			const maxYear = max ? Math.trunc(max / 10000) : 0;
			if (minYear && year < minYear) {
				yearInp.value = String(viewYear);
				return;
			}
			if (maxYear && year > maxYear) {
				yearInp.value = String(viewYear);
				return;
			}
			viewYear = year;
			render();
		});
		yearInp.addEventListener(`keydown`, (ev) => {
			if (ev.key !== `Enter`) return;
			ev.preventDefault();
			yearInp.blur();
		});

		popup.addEventListener(`click`, (ev) => {
			const nav = ev.target.closest(`button.qbuy-nav`);
			if (nav instanceof HTMLButtonElement) {
				shiftMonth(Number.parseInt(nav.dataset.step || `0`, 10));
				render();
				return;
			}
			const dayBtn = ev.target.closest(`button[data-day]`);
			if (!(dayBtn instanceof HTMLButtonElement)) return;
			if (!hidden || !trigger || dayBtn.disabled) return;
			const day = Number.parseInt(dayBtn.dataset.day || `0`, 10);
			if (day < 1 || day > 31) return;
			selected = ymd(viewYear, viewMonth, day);
			if (!inRange(selected)) return;
			hidden.value = as8(selected);
			trigger.value = asDMY(selected);
			hidden.dispatchEvent(new Event(`input`, { bubbles: true }));
			hidden.dispatchEvent(new Event(`change`, { bubbles: true }));
			closePopup();
		});
	};

	const openPopup = (nextWrap) => {
		ensurePopup();
		const h = nextWrap.querySelector(`input[data-buy-hidden="1"]`);
		const t = nextWrap.querySelector(`input[data-buy-trigger="1"]`);
		if (!(h instanceof HTMLInputElement) || !(t instanceof HTMLInputElement)) return;

		wrap = nextWrap;
		hidden = h;
		trigger = t;
		min = parse8(h.dataset.min || ``);
		max = parse8(h.dataset.max || ``);
		if (min && max && min > max) max = min;

		selected = parse8(h.value);
		if (!selected) selected = min || today();
		if (min && selected < min) selected = min;
		if (max && selected > max) selected = max;
		viewYear = Math.trunc(selected / 10000);
		viewMonth = Math.trunc((selected % 10000) / 100);

		render();
		setPopupPos();
		popup.classList.add(`is-open`);
	};

	document.addEventListener(`click`, (ev) => {
		const target = ev.target;
		if (!(target instanceof HTMLElement)) return;
		const triggerWrap = target.closest(`[data-buy="1"] .qbuy-trigger-wrap`);
		if (triggerWrap instanceof HTMLElement) {
			ev.preventDefault();
			const nextWrap = triggerWrap.closest(`[data-buy="1"]`);
			if (!(nextWrap instanceof HTMLElement)) return;
			if (wrap === nextWrap && popup && popup.classList.contains(`is-open`)) { closePopup(); return }
			openPopup(nextWrap);
			return;
		}
		if (!popup || !popup.classList.contains(`is-open`)) return;
		if (popup.contains(target)) return;
		if (wrap && wrap.contains(target)) return;
		closePopup();
	});

	document.addEventListener(`keydown`, (ev) => {
		if (ev.key === `Escape`) closePopup();
	});

	window.addEventListener(`resize`, () => {
		if (popup && popup.classList.contains(`is-open`)) setPopupPos();
	});
	window.addEventListener(`scroll`, () => {
		if (popup && popup.classList.contains(`is-open`)) setPopupPos();
	}, true);
})();
