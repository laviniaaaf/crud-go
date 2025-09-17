const apiUrl = '/itens';

const form = document.getElementById('form-item');
const itemIdInput = document.getElementById('item-id');
const nomeInput = document.getElementById('nome');
const precoInput = document.getElementById('preco');
const submitButton = document.getElementById('btn-submit');
const cancelButton = document.getElementById('btn-cancelar');

//  carregar os itens da API
async function carregarItens() {
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
    console.error('Erro ao carregar itens:', err);
    }
}

// CREATE / UPDATE: 
form.addEventListener('submit', async (event) => {
    event.preventDefault();

    const id = itemIdInput.value;
    const item = {
    nome: nomeInput.value,
    preco: parseFloat(precoInput.value)
    };

    if (id) {
    await atualizarItem(id, item);
    } else {
    await adicionarItem(item);
    }
});

async function adicionarItem(item) {
    try {
    const res = await fetch(apiUrl, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(item)
    });
    if (!res.ok) {
        const msg = await res.text();
        alert('Erro ao adicionar item: ' + msg);
        return;
    }
    form.reset();
    carregarItens();
    } catch (err) {
    console.error('Erro ao adicionar item:', err);
    alert('Erro ao adicionar item (network): ' + err.message);
    }
}

function prepararEdicao(item) {
    itemIdInput.value = item.id;
    nomeInput.value = item.nome;
    precoInput.value = item.preco;

    submitButton.textContent = 'Atualizar Item';
    cancelButton.style.display = 'inline-block';
    window.scrollTo(0, 0);
}

async function atualizarItem(id, item) {
    try {
    const res = await fetch(`${apiUrl}/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(item)
    });
    if (!res.ok) {
        const msg = await res.text();
        alert('Erro ao atualizar item: ' + msg);
        return;
    }
    cancelarEdicao();
    carregarItens();
    } catch (err) {
    console.error('Erro ao atualizar item:', err);
    alert('Erro ao atualizar item (network): ' + err.message);
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
        <button class="btn-editar" onclick='prepararEdicao(${JSON.stringify(item)})'>Editar</button>
        <button class="btn-excluir" onclick="excluirItem(${item.id})">Excluir</button>
        </div> 
    `;
    ulResultado.appendChild(li);
    } catch (err) {
    console.error('Erro ao buscar item:', err);
    ulResultado.innerHTML = `<li style="color:red;">Erro ao buscar item</li>`;
    }
}

function cancelarEdicao() {
    form.reset();
    itemIdInput.value = '';
    submitButton.textContent = 'Adicionar Item';
    cancelButton.style.display = 'none';
}

async function excluirItem(id) {
    if (!confirm('Você certeza que deseja excluir este item???')) return;
    try {
    const res = await fetch(`${apiUrl}/${id}`, { method: 'DELETE' });
    if (!res.ok) {
        const msg = await res.text();
        alert('Erro ao excluir item: ' + msg);
        return;
    }
    carregarItens();
    } catch (err) {
    console.error('Erro ao excluir item:', err);
    alert('Erro ao excluir item (network): ' + err.message);
    }
}

carregarItens();