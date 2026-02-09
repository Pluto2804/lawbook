// Main JavaScript for Lawbook

document.addEventListener('DOMContentLoaded', function() {
    // Flash message auto-dismiss
    const flashMessage = document.querySelector('.flash');
    if (flashMessage) {
        setTimeout(() => {
            flashMessage.style.opacity = '0';
            setTimeout(() => flashMessage.remove(), 300);
        }, 5000);
    }
    
    // Form validation enhancement
    const forms = document.querySelectorAll('form[novalidate]');
    forms.forEach(form => {
        form.addEventListener('submit', function(e) {
            // Custom validation can be added here
        });
    });
});
