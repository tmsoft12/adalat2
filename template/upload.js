// Kayıt Formu İşlemleri
document.getElementById('registerForm').addEventListener('submit', async (event) => {
    event.preventDefault();

    const username = document.getElementById('registerUsername').value;
    const password = document.getElementById('registerPassword').value;

    const response = await fetch('http://192.168.0.111/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            username: username,
            password: password,
        }),
    });

    const data = await response.json();
    document.getElementById('registerMessage').innerText = data.message;
});

// Giriş Yapma Formu İşlemleri
document.getElementById('loginForm').addEventListener('submit', async (event) => {
    event.preventDefault();

    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;

    const response = await fetch('http://192.168.0.111:5000/api/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            username: username,
            password: password,
        }),
        credentials: 'include', // Cookie'leri dahil etmek için kullanılır
    });

    const data = await response.json();
    document.getElementById('loginMessage').innerText = data.message;
});

// Korunan Route'a Erişim
document.getElementById('accessProtectedBtn').addEventListener('click', async () => {
    const response = await fetch('http://192.168.0.111:5000/api/auth/protected', {
        method: 'GET',
        credentials: 'include',  // Cookie'leri otomatik olarak ekler
    });

    if (response.status === 200) {
        const data = await response.json();
        document.getElementById('protectedMessage').innerText = `Protected Route Accessed. User ID: ${data.userID}`;
    } else {
        document.getElementById('protectedMessage').innerText = 'Access denied';
    }
});

// Token Yenileme
document.getElementById('refreshTokenBtn').addEventListener('click', async () => {
    const response = await fetch('http://192.168.0.111:5000/api/auth/refresh', {
        method: 'POST',
        credentials: 'include', // Cookie'leri dahil etmek için kullanılır
    });

    const data = await response.json();
    document.getElementById('refreshMessage').innerText = data.message;
});
