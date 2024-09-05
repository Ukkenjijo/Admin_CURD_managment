$(document).ready(function () {
    $('[data-toggle="tooltip"]').tooltip();
});

document.addEventListener("DOMContentLoaded", function () {
    document.querySelectorAll(".edit").forEach(function (button) {
        button.addEventListener("click", function (event) {
            event.preventDefault();

            // Get user data from data attributes
            const userId = this.getAttribute("data-id");
            const username = this.getAttribute("data-username");
            const email = this.getAttribute("data-email");

            // Populate the modal with user data
            document.getElementById("editUserId").value = userId;
            document.getElementById("editUsername").value = username;
            document.getElementById("editEmail").value = email;

            // Show the modal
            $('#editUserModal').modal('show');
        });
    });

    // Handle form submission
    document.getElementById("editUserForm").addEventListener("submit", function (event) {
        event.preventDefault();

        // Get form data
        const formData = new FormData(this);

        // Send the update request via Fetch API
        fetch(`/admin/edit/${formData.get('id')}`, {
            method: 'POST', // Change to 'PUT' or 'PATCH' as necessary
            body: formData
        }).then(response => {
            if (response.ok) {
                alert("User updated successfully");
                window.location.reload(); // Reload the page to reflect changes
            } else {
                alert("Failed to update user");
            }
        });
    });
    // New: Handle Add User form submission
    document.getElementById("createUserForm").addEventListener("submit", function (event) {
        event.preventDefault();  // Prevent default form submission

        const formData = new FormData(this);

        fetch("/admin/create", {
            method: "POST",  // Adjust the method if necessary
            body: formData
        }).then(response => {
            if (response.ok) {
                alert("User created successfully");
                window.location.reload();  // Reload the page to reflect changes
            } else {
                alert("Failed to create user");
            }
        });
    });
});
