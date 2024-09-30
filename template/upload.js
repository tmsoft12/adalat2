// Function to handle form submission
async function handleFormSubmit(event, url, formId) {
    event.preventDefault(); // Prevent the default form submission

    const form = document.getElementById(formId);
    const formData = new FormData(form); // Collect the form data

    try {
        const response = await fetch(url, {
            method: 'POST',
            body: formData
        });

        if (response.ok) {
            const result = await response.json();
            alert(`Success: ${JSON.stringify(result)}`);
        } else {
            alert(`Error: ${response.statusText}`);
        }
    } catch (error) {
        alert(`Request failed: ${error.message}`);
    }
}

// Add event listeners to forms
document.getElementById('newsForm').addEventListener('submit', function(event) {
    handleFormSubmit(event, 'http://localhost:5000/api/admin/news', 'newsForm');
});

document.getElementById('bannerForm').addEventListener('submit', function(event) {
    handleFormSubmit(event, 'http://localhost:5000/api/admin/banner', 'bannerForm');
});

document.getElementById('employerForm').addEventListener('submit', function(event) {
    handleFormSubmit(event, 'http://localhost:5000/api/admin/employer', 'employerForm');
});

document.getElementById('mediaForm').addEventListener('submit', function(event) {
    handleFormSubmit(event, 'http://localhost:5000/api/admin/media', 'mediaForm');
});
