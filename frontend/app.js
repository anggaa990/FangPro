const API_BASE = 'http://localhost:8080';

// Set default datetime to now
document.addEventListener('DOMContentLoaded', function() {
    const now = new Date();
    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
    document.getElementById('priceDate').value = now.toISOString().slice(0, 16);
    
    // Load initial data
    fetchWeather();
    fetchPrices();
});

// Auto-Fetch Prices from Scraping
async function autoFetchPrices() {
    const btn = document.getElementById('autoFetchBtn');
    const originalText = btn.innerHTML;
    
    btn.disabled = true;
    btn.innerHTML = '⏳ Mengambil data...';
    
    try {
        const response = await fetch(`${API_BASE}/harga/fetch`, {
            method: 'POST'
        });
        
        const data = await response.json();
        
        if (response.ok) {
            alert('✅ ' + data.message);
            // Refresh prices table
            fetchPrices();
        } else {
            throw new Error(data.message || 'Gagal fetch harga');
        }
    } catch (error) {
        console.error('Error auto-fetching prices:', error);
        alert('❌ Gagal mengambil data harga: ' + error.message);
    } finally {
        btn.disabled = false;
        btn.innerHTML = originalText;
    }
}

// Fetch Weather Data
async function fetchWeather() {
    const region = document.getElementById('regionSelect').value;
    const weatherLoading = document.getElementById('weatherLoading');
    const weatherData = document.getElementById('weatherData');
    
    weatherLoading.style.display = 'block';
    weatherData.style.display = 'none';
    
    try {
        const response = await fetch(`${API_BASE}/cuaca?region=${region}`);
        const data = await response.json();
        
        document.getElementById('temp').textContent = data.temp.toFixed(1);
        document.getElementById('humidity').textContent = data.humidity;
        document.getElementById('rain').textContent = data.rain_mm.toFixed(1);
        
        weatherLoading.style.display = 'none';
        weatherData.style.display = 'block';
    } catch (error) {
        console.error('Error fetching weather:', error);
        weatherLoading.textContent = '❌ Gagal memuat data cuaca';
    }
}

// Fetch Advanced Recommendation
async function fetchAdvancedRecommendation() {
    const region = document.getElementById('regionSelect').value;
    const loading = document.getElementById('recommendationLoading');
    const dataDiv = document.getElementById('recommendationData');
    
    loading.style.display = 'block';
    dataDiv.style.display = 'none';
    
    try {
        const response = await fetch(`${API_BASE}/rekomendasi/advanced?region=${region}`);
        const data = await response.json();
        
        console.log('Recommendation data:', data);
        
        // Status Badge
        const statusBadge = document.getElementById('statusBadge');
        statusBadge.textContent = data.main_advice;
        statusBadge.className = `status-badge status-${data.status}`;
        
        // Main Advice
        document.getElementById('mainAdvice').textContent = data.main_advice;
        
        // Detailed Advice Cards
        document.getElementById('plantingAdvice').textContent = data.planting_advice;
        document.getElementById('harvestAdvice').textContent = data.harvest_advice;
        document.getElementById('dryingAdvice').textContent = data.drying_advice;
        document.getElementById('irrigationAdvice').textContent = data.irrigation_advice;
        document.getElementById('pestWarning').textContent = data.pest_warning;
        
        loading.style.display = 'none';
        dataDiv.style.display = 'block';
    } catch (error) {
        console.error('Error fetching recommendation:', error);
        loading.textContent = '❌ Gagal memuat rekomendasi';
    }
}

// Fetch Recommendation (simple - backward compatible)
async function fetchRecommendation() {
    const loading = document.getElementById('recommendationLoading');
    const dataDiv = document.getElementById('recommendationData');
    const textDiv = document.getElementById('recommendationText');
    
    loading.style.display = 'block';
    dataDiv.style.display = 'none';
    
    try {
        const response = await fetch(`${API_BASE}/rekomendasi`);
        const data = await response.json();
        
        textDiv.textContent = data.recommendation;
        
        // Change color based on recommendation
        if (data.recommendation.includes('cocok untuk tanam')) {
            textDiv.className = 'recommendation-box';
        } else if (data.recommendation.includes('cocok untuk pengeringan')) {
            textDiv.className = 'recommendation-box';
        } else {
            textDiv.className = 'recommendation-box neutral';
        }
        
        loading.style.display = 'none';
        dataDiv.style.display = 'block';
    } catch (error) {
        console.error('Error fetching recommendation:', error);
        loading.textContent = '❌ Gagal memuat rekomendasi';
    }
}

// Add Price
async function addPrice(event) {
    event.preventDefault();
    
    const region = document.getElementById('priceRegion').value;
    const price = parseFloat(document.getElementById('priceValue').value);
    const unit = document.getElementById('priceUnit').value;
    const source = document.getElementById('priceSource').value;
    const recordedAt = document.getElementById('priceDate').value;
    
    const alertDiv = document.getElementById('formAlert');
    
    try {
        const response = await fetch(`${API_BASE}/harga/add`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                region: region,
                price: price,
                unit: unit,
                source: source,
                recorded_at: recordedAt
            })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            alertDiv.className = 'alert alert-success';
            alertDiv.textContent = '✅ ' + data.message;
            alertDiv.style.display = 'block';
            
            // Reset form
            document.getElementById('priceForm').reset();
            
            // Set datetime to now again
            const now = new Date();
            now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
            document.getElementById('priceDate').value = now.toISOString().slice(0, 16);
            
            // Refresh prices table
            fetchPrices();
            
            // Hide alert after 3 seconds
            setTimeout(() => {
                alertDiv.style.display = 'none';
            }, 3000);
        } else {
            throw new Error(data.message || 'Gagal menambahkan data');
        }
    } catch (error) {
        console.error('Error adding price:', error);
        alertDiv.className = 'alert alert-error';
        alertDiv.textContent = '❌ ' + error.message;
        alertDiv.style.display = 'block';
    }
}

// Fetch Prices
async function fetchPrices() {
    const loading = document.getElementById('pricesLoading');
    const dataDiv = document.getElementById('pricesData');
    
    loading.style.display = 'block';
    dataDiv.innerHTML = '';
    
    try {
        const response = await fetch(`${API_BASE}/harga`);
        
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        
        const data = await response.json();
        
        console.log('Fetched prices:', data); // Debug log
        
        if (data && Array.isArray(data) && data.length > 0) {
            const table = document.createElement('div');
            table.className = 'table-container';
            table.innerHTML = `
                <table>
                    <thead>
                        <tr>
                            <th>No</th>
                            <th>Region</th>
                            <th>Harga</th>
                            <th>Satuan</th>
                            <th>Sumber</th>
                            <th>Tanggal Catat</th>
                            <th>Input Sistem</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${data.map((item, index) => `
                            <tr>
                                <td>${index + 1}</td>
                                <td><span class="badge badge-success">${item.region}</span></td>
                                <td><strong>Rp ${formatNumber(item.price)}</strong></td>
                                <td>${item.unit}</td>
                                <td>${item.source}</td>
                                <td>${formatDate(item.recorded_at)}</td>
                                <td>${formatDate(item.created_at)}</td>
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            `;
            dataDiv.appendChild(table);
        } else {
            dataDiv.innerHTML = `
                <div class="empty-state">
                    <div class="empty-state-icon">📦</div>
                    <p>Belum ada data harga tembakau. Klik "Auto-Fetch Harga" untuk mengambil data.</p>
                </div>
            `;
        }
        
        loading.style.display = 'none';
    } catch (error) {
        console.error('Error fetching prices:', error);
        loading.style.display = 'none';
        dataDiv.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">❌</div>
                <p>Gagal memuat data harga: ${error.message}</p>
            </div>
        `;
    }
}

// Helper: Format Number
function formatNumber(num) {
    return new Intl.NumberFormat('id-ID').format(num);
}

// Helper: Format Date
function formatDate(dateString) {
    if (!dateString) return '-';
    
    const date = new Date(dateString);
    const options = { 
        year: 'numeric', 
        month: 'short', 
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    };
    
    return date.toLocaleDateString('id-ID', options);
}