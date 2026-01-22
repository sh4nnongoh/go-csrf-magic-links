// Read from meta tag and store in localStorage
const csrfMeta = document.querySelector('meta[name="csrf-token"]')?.content;
if (csrfMeta) {
  localStorage.setItem('csrfToken', csrfMeta);
}

// Attach CSRF token from meta tag or localstorage to HTMX request
document.addEventListener('htmx:configRequest', (event) => {
  event.detail.headers['X-CSRF-Token'] =
    csrfMeta || localStorage.getItem('csrfToken')
});
