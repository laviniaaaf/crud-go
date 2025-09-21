const apiUrl = '/itens';

const form = document.getElementById('form-item');
const itemIdInput = document.getElementById('item-id');
const nomeInput = document.getElementById('nome');
const precoInput = document.getElementById('preco');
const submitButton = document.getElementById('btn-submit');
const cancelButton = document.getElementById('btn-cancelar');

// load items from the API 
async function loadItems() {
    try {
    const res = await fetch(apiUrl);
    const itens = await res.json();

    const ul = document.getElementById('lista-itens');
    ul.innerHTML = '';

    if (!itens || itens.length === 0) {
        ul.innerHTML = '<li>Nenhum item cadastrado</li>';
        return;
    }

    itens.forEach(item => {
        const li = document.createElement('li');
        li.innerHTML = `
        <span><b>ID:</b> ${item.id} — ${item.nome}: R$ ${item.preco.toFixed(2)}</span>
        <div class="item-actions">
            <button class="btn-editar" onclick='prepararEdicao(${JSON.stringify(item)})'>Editar</button>
            <button class="btn-excluir" onclick="excluirItem(${item.id})">Excluir</button>
        </div>
        `;
        ul.appendChild(li);
    });
    } catch (err) {
    console.error('Error loading items:', err);
    }
}


form.addEventListener('submit', async (event) => {
    event.preventDefault();

    const id = itemIdInput.value;
    const item = {
    nome: nomeInput.value,
    preco: parseFloat(precoInput.value)
    };

    if (id) {
    await updateItem(id, item);
    } else {
    await addItem(item);
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
        alert('Error adding item: ' + msg);
        return;
    }
    form.reset();
    loadItems();
    } catch (err) {
    console.error('Error adding item:', err);
    alert('Error adding item (network): ' + err.message);
    }
}

function prepareEdit(item) {
    itemIdInput.value = item.id;
    nomeInput.value = item.nome;
    precoInput.value = item.preco;

    submitButton.textContent = 'Atualizar Item';
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
        alert('Error updating item: ' + msg);
        return;
    }
    cancelEdit();
    loadItems();
    } catch (err) {
    console.error('Error updating item:', err);
    alert('Error updating item (network): ' + err.message);
    }
}

async function buscarItem() {

    const id = document.getElementById('busca').value;
    const ulResultado = document.getElementById('resultado-busca');
    ulResultado.innerHTML = '';

    if (!id) {
    alert("Informe um ID para buscar.");
    return;
    }

    try {
    const res = await fetch(`${apiUrl}/${id}`);

    if (!res.ok) {
        ulResultado.innerHTML = `<li style="color:red;">Item não encontrado</li>`;
        return;
    }

    const item = await res.json();

    const li = document.createElement('li');
    li.innerHTML = `
        <span><b>ID:</b> ${item.id} — ${item.nome}: R$ ${item.preco.toFixed(2)}</span>
        <div class="item-actions">
        <button class="btn-editar" onclick='prepareEdit(${JSON.stringify(item)})'>Edit</button>
        <button class="btn-excluir" onclick="deleteItem(${item.id})">Delete</button>
        </div> 
    `;
    ulResultado.appendChild(li);
    } catch (err) {
    console.error('Error searching item:', err);
    ulResultado.innerHTML = `<li style="color:red;">Error searching item</li>`;
    }
}

function cancelEdit() {
    form.reset();
    itemIdInput.value = '';
    submitButton.textContent = 'Add Item';
    cancelButton.style.display = 'none';
}

async function deleteItem(id) {
    if (!confirm('You sure you want to delete this item???')) return;
    try {
    const res = await fetch(`${apiUrl}/${id}`, { method: 'DELETE' });
    if (!res.ok) {
        const msg = await res.text();
        alert('Error deleting item: ' + msg);
        return;
    }
    loadItems();
    } catch (err) {
    console.error('Error deleting item:', err);
    alert('Error deleting item (network): ' + err.message);
    }
}

loadItems();