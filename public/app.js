let currencies = [];

async function loadCurrencies() {
    const response = await fetch('/codes');
    currencies = await response.json();

    populateSelect(document.getElementById('from'));
    populateSelect(document.getElementById('to'));

    document.querySelectorAll('.fav-from, .fav-to').forEach(populateSelect);

    loadFavoritesIntoDropdowns();
}

function populateSelect(select) {
    for (const currency of currencies) {
        const option = document.createElement('option');
        option.value = currency.Code;
        option.textContent = `${currency.Code} - ${currency.Name}`;
        select.appendChild(option);
    }
}

loadCurrencies();

let debounceTimer;

function scheduleConvert() {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(convert, 400);
}

async function convert() {
    const amount = document.getElementById('amount').value;
    const from = document.getElementById('from').value;
    const to = document.getElementById('to').value;
    const resultDiv = document.getElementById('result');

    if (!amount || !from || !to) {
        resultDiv.textContent = '';
        return;
    }

    const response = await fetch(`/convert?amount=${amount}&from=${from}&to=${to}`);

    if (!response.ok) {
        resultDiv.textContent = 'Conversion failed';
        return;
    }

    const data = await response.json();
    const result = data.Result.toFixed(2);
    resultDiv.textContent = `${amount} ${from} = ${result} ${to}`;

    saveToHistory({ amount, from, to, result });
}

document.getElementById('amount').addEventListener('input', scheduleConvert);
document.getElementById('from').addEventListener('change', scheduleConvert);
document.getElementById('to').addEventListener('change', scheduleConvert);

function loadHistory() {
    const raw = localStorage.getItem('history');
    return raw ? JSON.parse(raw) : [];
}

function saveToHistory(entry) {
    const history = loadHistory();
    history.unshift(entry); // newest first
    if (history.length > 5) {
        history.pop(); // drop oldest
    }
    localStorage.setItem('history', JSON.stringify(history));
    renderHistory();
}

function renderHistory() {
    const history = loadHistory();
    const historyList = document.getElementById('history-list');
    historyList.innerHTML = '';

    for (const entry of history) {
        const item = document.createElement('div');
        item.textContent = `${entry.amount} ${entry.from} = ${entry.result} ${entry.to}`;
        historyList.appendChild(item);
    }
}

renderHistory();

function loadFavorites() {
    const raw = localStorage.getItem('favorites');
    return raw ? JSON.parse(raw) : [null, null, null];
}

function saveFavorites(favorites) {
    localStorage.setItem('favorites', JSON.stringify(favorites));
}

function loadFavoritesIntoDropdowns() {
    const favorites = loadFavorites();

    document.querySelectorAll('.fav-from').forEach(select => {
        const index = select.dataset.index;
        if (favorites[index]) select.value = favorites[index].from;
    });

    document.querySelectorAll('.fav-to').forEach(select => {
        const index = select.dataset.index;
        if (favorites[index]) select.value = favorites[index].to;
    });
}

function saveFavoriteRow(index) {
    const fromSelect = document.querySelector(`.fav-from[data-index="${index}"]`);
    const toSelect = document.querySelector(`.fav-to[data-index="${index}"]`);

    const favorites = loadFavorites();
    favorites[index] = { from: fromSelect.value, to: toSelect.value };
    saveFavorites(favorites);
}

document.querySelectorAll('.fav-from, .fav-to').forEach(select => {
    select.addEventListener('change', () => saveFavoriteRow(select.dataset.index));
});

document.querySelectorAll('.fav-load').forEach(button => {
    button.addEventListener('click', () => {
        const index = button.dataset.index;
        const favorites = loadFavorites();
        const pair = favorites[index];
        if (!pair) return;

        document.getElementById('from').value = pair.from;
        document.getElementById('to').value = pair.to;
        convert();
    });
});