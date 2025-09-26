const apiUrl = '/bills';

const form = document.getElementById('form-item');
const itemIdInput = document.getElementById('item-id');
const embasaInput = document.getElementById('embasa');
const coelbaInput = document.getElementById('coelba');
const createdAtInput = document.getElementById('created-at');
const updatedAtInput = document.getElementById('updated-at');
const submitButton = document.getElementById('btn-submit');
const cancelButton = document.getElementById('btn-cancel');

function nowISO() {
    return new Date().toISOString().slice(0, 16);
}

createdAtInput.value = nowISO();

form.addEventListener('submit', async (event) => {
    event.preventDefault();

    const id = itemIdInput.value;
    const bill = {
        embasa: embasaInput.value,
        coelba: coelbaInput.value,
        created_at: createdAtInput.value || nowISO(),
        updated_at: id ? nowISO() : null  // add updated_at if the person edited the item
    };

    if (id) {
        await updateItem(id, bill);
    } else {
        await addItem(bill);
    }
});

async function loadItems() {
    try {
        const res = await fetch(apiUrl);
        const items = await res.json();

        const ul = document.getElementById('item-list');
        ul.innerHTML = '';

        if (!items || items.length === 0) {
            ul.innerHTML = '<li>No items registered</li>';
        }

        
        items.forEach(item => {
            const li = document.createElement('li');
            li.classList.add("bill-card");
            li.innerHTML = `
                <div class="bill-header">
                    <span class="bill-id">${item.id}</span>
                </div>
                <div class="bill-body">
                    <p><b>EMBASA:</b> R$ ${parseFloat(item.embasa).toFixed(2)}</p>
                    <p><b>COELBA:</b> R$ ${parseFloat(item.coelba).toFixed(2)}</p>
                    <p><b>Created:</b> ${formatDate(item.created_at)}</p>
                    <p><b>Updated:</b> ${formatDate(item.updated_at)}</p>
                </div>
                <div class="bill-actions">
                    <button class="btn-edit" onclick='prepareEdit(${JSON.stringify(item)})'>Edit</button>
                    <button class="btn-delete" onclick="deleteItem('${item.id}')">Delete</button>
                </div>
            `;
            ul.appendChild(li);
        });
        
    } catch (err) {
        console.error('Error loading items:', err);
    }
}

// convert datetime-local to RFC3339
function toRFC3339(datetimeLocalValue) {
    if (!datetimeLocalValue) return new Date().toISOString();
    return new Date(datetimeLocalValue).toISOString();
}

form.addEventListener('submit', async (event) => {
    event.preventDefault();

    const id = itemIdInput.value;

    // valid numbers
    const embasaValue = parseFloat(embasaInput.value) || 0;
    const coelbaValue = parseFloat(coelbaInput.value) || 0;

    const bill = {
        embasa: embasaValue,
        coelba: coelbaValue,
        created_at: toRFC3339(createdAtInput.value)
    };

    if (id) {
        bill.updated_at = new Date().toISOString();
        console.log("PUT body:", bill); 
        await updateItem(id, bill);
    } else {
        console.log("POST body:", bill); 
        await addItem(bill);
    }
});


async function addItem(item) {
    try {
        const res = await fetch(apiUrl, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(item)
        });
        if (!res.ok) {
            const msg = await res.text();
            return;
        }
        form.reset();
        loadItems();
    } catch (err) {
        
    }
}

function formatDate(dateString) {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleString('pt-BR');
}

function prepareEdit(item) {
    itemIdInput.value = item.id;
    embasaInput.value = parseFloat(item.embasa).toFixed(2);
    coelbaInput.value = parseFloat(item.coelba).toFixed(2);
    createdAtInput.value = item.created_at ? item.created_at.slice(0, 16) : nowISO();

    const updatedWrapper = document.getElementById('updated-at-wrapper');
    updatedWrapper.style.display = 'block';
    updatedAtInput.value = nowISO();

    submitButton.textContent = 'Update';
    cancelButton.style.display = 'inline-block';
    window.scrollTo(0, 0);
}

async function updateItem(id, item) {
    try {
        const res = await fetch(`${apiUrl}/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(item)
        });
        if (!res.ok) {
            const msg = await res.text();
            return;
        }
        cancelEdit();
        loadItems();
    } catch (err) {
       
    }
}

async function searchItem() {
    const id = document.getElementById('search').value;
    const ulResult = document.getElementById('search-result');
    ulResult.innerHTML = '';

    if (!id) {
        alert("Enter an ID to search.");
        return;
    }

    try {
        const res = await fetch(`${apiUrl}/${id}`);

        if (!res.ok) {
            ulResult.innerHTML = `<li style="color:red;">Item not found</li>`;
            return;
        }

        const item = await res.json();

        const li = document.createElement('li');
        li.innerHTML = `
            <span>
                <b>ID:</b> ${item.id} | 
                <b>EMBASA:</b> R$ ${parseFloat(item.embasa).toFixed(2)} | 
                <b>COELBA:</b> R$ ${parseFloat(item.coelba).toFixed(2)} | 
                <b>Created:</b> ${formatDate(item.created_at)} | 
                <b>Updated:</b> ${formatDate(item.updated_at)}
            </span>
            <div class="item-actions">
                <button class="btn-edit" onclick='prepareEdit(${JSON.stringify(item)})'>Edit</button>
                <button class="btn-delete" onclick="deleteItem('${item.id}')">Delete</button>
            </div> 
        `;
        ulResult.appendChild(li);
    } catch (err) {
        console.error('Error searching item:', err);
        ulResult.innerHTML = `<li style="color:red;">Error searching item</li>`;
    }
}

function cancelEdit() {
    form.reset();
    itemIdInput.value = '';
    createdAtInput.value = nowISO();

    document.getElementById('updated-at-wrapper').style.display = 'none';
    updatedAtInput.value = '';

    submitButton.textContent = 'Save';
    cancelButton.style.display = 'none';
}

async function deleteItem(id) {
    if (!confirm('Are you sure you want to delete this item???')) return;
    try {
        await fetch(`${apiUrl}/${id}`, { method: 'DELETE' });
        loadItems();
    } catch (err) {

    }
}

loadItems();
