function printTable() {
    // Сохраняем HTML таблицы
    const tableHTML = document.getElementById('namazTable').outerHTML;
    
    // Создаем новое окно для печати
    const printWindow = window.open('', '_blank');
    printWindow.document.write(`
        <!DOCTYPE html>
        <html>
        <head>
            <title>Расписание намазов (г. Норильск)</title>
            <style>
                body { font-family: Arial, sans-serif; margin: 20px; }
                table { width: 100%; border-collapse: collapse; }
                th, td { border: 1px solid #ddd; padding: 8px; text-align: center; }
                th { background-color: #f2f2f2; font-weight: bold; }
                @media print { body { margin: 0; } }
            </style>
        </head>
        <body>
            <h1 style="text-align: center; margin-bottom: 5px;">Расписание намазов</h1>
            <h3 style="text-align: center; margin-bottom: 20px;">г. Норильск</h3>
            ${tableHTML}
        </body>
        </html>
    `);
    
    printWindow.document.close();
    printWindow.focus();
    printWindow.print();
    printWindow.close();
}