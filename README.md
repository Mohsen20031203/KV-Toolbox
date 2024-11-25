# ManageDB

ManageDB is a simple and user-friendly application for managing key-value databases. It allows users to create, edit, and delete keys, browse database contents, and perform various operations effortlessly.

![ManageDB Logo](cmd/KV-Toolbox/icon.png)

---

## Features

- Create new key-value databases.
- Add, edit, or delete keys easily.
- Browse database contents with an intuitive interface.
- Cross-platform support (details below).
- Lightweight and efficient, built with [Fyne](https://fyne.io/).

---

## Installation

### macOS
1. Click on **Install** below to download the program:
   - [INSTALL](https://github.com/Mohsen20031203/KV-Toolbox/releases/download/v1.0.0/ManageDB.app.zip)

2. When you try to open the app, macOS might display the following error:  
   _“This application ‘ManageDB.app’ can't be opened”._

   Follow these steps to resolve the issue:

    #### Step 1: Grant executable permissions
   Open the **terminal**, first go to the location of your application and run the following command to make the application executable:

   ```bash
   chmod +x ManageDB.app/Contents/MacOS/'KV-Toolbox'
#### Step 2: Allow the App in Security Settings

1. Open **System Preferences**.
2. Go to **Security & Privacy** > **General** tab.
3. Look for a message stating:  
   _“‘ManageDB.app’ was blocked because it is not from an identified developer.”_
4. Click **Allow Anyway** or **Open Anyway**.

#### Step 3: Reopen the App

1. Try opening `ManageDB.app` again.
2. The app should now run without any issues.



### Windows
1. Download the `ManageDB.exe` file from the [releases section](#).
2. Double-click to run the installer and follow the setup instructions.
3. Launch the app from your desktop or start menu.

### Linux
1. Download the `ManageDB.AppImage` file from the [releases section](#).
2. Make the file executable using: `chmod +x ManageDB.AppImage`.
3. Run the file using: `./ManageDB.AppImage`.

---

## How to Use

1. Open the application.
2. **Create a new database**: To create a new database, click on the "+" button and select the database you want, then click on "Create Database" in the opened window and enter the address you want.
2. **Open a new database**: To open a database, click the "+" button and select the database you want, then click the "Open folder" button in the opened window and find your database.
3. **Add keys**: Use the "Add" option to create new key-value pairs.
    - You can also assign files, images, or other types of data as the value for a key. This feature allows you to efficiently manage and store additional resources, such as images and files, directly within your database.
4. **Edit keys**: Find the key you want and click on it to change its value.
5. **Delete keys**: Click the "delete" key in the main window and then enter the key you want.
6. **Search keys**: Click on the search button and enter the key you want.

---

## System Requirements

- **macOS**: Version 10.14 or later.
- **Windows**: Version 7 or later.
- **Linux**: Kernel version 4.0 or later.
- **Storage**: At least 50 MB of free space.

---

## Screenshots

Below are some screenshots showcasing the application in action:

1. **Main Interface**  
   ![Main Interface](./screenshots/main_interface.png)  
   _A simple, clean layout to manage your databases._

2. **Editing Keys**  
   ![Editing Keys](./screenshots/edit_keys.png)  
   _Easily edit existing keys and values._

3. **Database Browser**  
   ![Database Browser](./screenshots/database_browser.png)  
   _Quickly browse through your data._

---