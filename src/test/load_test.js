import http from 'k6/http';
import { check, sleep } from 'k6';
import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.1/index.js';


function generateRandomUsername() {
    const prefix = 'user';
    const randomNumber = Math.floor(10000 + Math.random() * 90000); // 5-значное число
    return `${prefix}${randomNumber}`;
}

// Функция для получения токенов
function getTokens(vus) {
    let tokens = [];

    for (let i = 0; i < vus; i++) {
        let username = generateRandomUsername();
        let authPayload = JSON.stringify({
            username: username,
            password: '12345pass',
        });

        let authRes = http.post('http://localhost:8080/api/auth', authPayload, {
            headers: { 'Content-Type': 'application/json' },
        });

        check(authRes, {
            'auth response status is 200': (r) => r.status === 200,
        });

        if (authRes.status === 200) {
            tokens.push(authRes.json('token'));
        } else {
            console.log(`Failed to get token for user ${username}`);
        }
    }

    return tokens;
}

// Функция setup вызывается перед запуском теста
export function setup() {
    let vus = 500; // Максимальное количество пользователей
    return getTokens(vus);
}

// Постепенное увеличение нагрузки
export let options = {
    stages: [
        { duration: '30s', target: 50 },  // Увеличение от 0 до 50 пользователей за 30 секунд
        { duration: '1m', target: 200 },  // Увеличение до 200 пользователей за 1 минуту
        { duration: '1m', target: 500 },  // Увеличение до 500 пользователей за 1 минуту
        { duration: '2m', target: 500 },  // Поддержание 500 пользователей в течение 2 минут
        { duration: '30s', target: 0 },   // Постепенное снижение до 0 пользователей за 30 секунд
    ],
    thresholds: {
        http_req_duration: ['p(99)<50'], // 99% запросов должны быть быстрее 50 мс
        http_req_failed: ['rate<0.0001'], // Ошибки меньше 0.01%
    },
};

// Основной тест
export default function (tokens) {
    let token = tokens[__VU - 1];

    let res = http.get('http://localhost:8080/api/info', {
        headers: {
            Authorization: token,
        },
    });

    check(res, {
        'status is 200': (r) => r.status === 200,
    });

    sleep(1);
}

// Функция для сохранения отчета
export function handleSummary(data) {
    return {
        'summary.txt': textSummary(data, { indent: ' ', enableColors: false }),
    };
}