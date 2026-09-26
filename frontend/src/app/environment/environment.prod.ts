// Ambiente de produção (usado em `ng build --configuration production`).
// apiUrl aponta para o backend publicado no Render.
// Precisa ser HTTPS (o frontend roda em HTTPS no Vercel), senão o navegador bloqueia
// por "Mixed Content".
export const environment = {
  production: true,
  apiUrl: 'https://dose-certa-backend.onrender.com',
};
