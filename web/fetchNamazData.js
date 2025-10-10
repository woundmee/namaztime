// Функция для загрузки данных
async function fetchNamazData(endpoint) {
    try {
      const response = await fetch(`http://localhost:8080${endpoint}`);
      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Ошибка загрузки данных:', error);
      return null;
    }
  }

  // Функция для отрисовки данных в таблице
  function renderNamazTable(data) {
    const tableBody = document.getElementById('namazTableBody');
    tableBody.innerHTML = ''; // Очищаем таблицу

    // Проверяем, является ли data массивом или объектом
    const dataToRender = Array.isArray(data) ? data : [data];

    dataToRender.forEach(item => {
      const row = document.createElement('tr');
      row.innerHTML = `
        <th scope="row">${item.Day}</th>
        <td>${item.Fajr}</td>
        <td>${item.Sunrise}</td>
        <td>${item.Zuhr}</td>
        <td>${item.Asr}</td>
        <td>${item.Magrib}</td>
        <td>${item.Isha}</td>
      `;
      tableBody.appendChild(row);
    });
  }

  // Обработчики кнопок
  const btnToday = document.getElementById('btnToday');
  const btnMonth = document.getElementById('btnMonth');

  btnToday.addEventListener('click', async () => {
    const data = await fetchNamazData('/namaztime/today');
    if (data) {
      renderNamazTable(data);
      btnToday.classList.replace('btn-outline-primary', 'btn-primary');
      btnMonth.classList.replace('btn-primary', 'btn-outline-primary');
    }
  });

  btnMonth.addEventListener('click', async () => {
    const data = await fetchNamazData('/namaztime/month');
    if (data) {
      renderNamazTable(data);
      btnMonth.classList.replace('btn-outline-primary', 'btn-primary');
      btnToday.classList.replace('btn-primary', 'btn-outline-primary');
    }
  });

  // Переключение темы
  const themeToggle = document.getElementById('themeToggle');
  const themeIcon = document.getElementById('themeIcon');
  const htmlElement = document.documentElement;

  themeToggle.addEventListener('click', () => {
    const currentTheme = htmlElement.getAttribute('data-bs-theme');
    
    if (currentTheme === 'light') {
      htmlElement.setAttribute('data-bs-theme', 'dark');
      themeIcon.classList.remove('bi-sun-fill');
      themeIcon.classList.add('bi-moon-fill');
    } else {
      htmlElement.setAttribute('data-bs-theme', 'light');
      themeIcon.classList.remove('bi-moon-fill');
      themeIcon.classList.add('bi-sun-fill');
    }
  });

  // Навигация между страницами
  const mainLink = document.getElementById('mainLink');
  const faqLink = document.getElementById('faqLink');
  const mainContent = document.getElementById('mainContent');
  const faqContent = document.getElementById('faqContent');

  mainLink.addEventListener('click', (e) => {
    e.preventDefault();
    mainContent.style.display = 'block';
    faqContent.style.display = 'none';
    mainLink.classList.add('active');
    faqLink.classList.remove('active');
  });

  faqLink.addEventListener('click', (e) => {
    e.preventDefault();
    mainContent.style.display = 'none';
    faqContent.style.display = 'block';
    faqLink.classList.add('active');
    mainLink.classList.remove('active');
  });

  // Автоматическая загрузка данных при загрузке страницы
  document.addEventListener('DOMContentLoaded', async () => {
    const data = await fetchNamazData('/namaztime/today');
    if (data) renderNamazTable(data);
  });