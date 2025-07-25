---
sidebar_position: 10
title: Notifications
description: How Notifications work
---

### **Intento Notification System: How It Works**

The Intento notification system is built for flexibility and speed. It plugs in to the chain's event and sends emails to users. As such, a notification system can also be operated by any other party. This notification system is run by team and is a complimentary service for integrators. Subscribing to flows is also part of the [submit page](./submit-page) and can be implemented directly into applications.

It starts with subscribing to a flow **via a Netlify function call**. To subscribe integrators send a request to:

```
https://portal.intento.zone/.netlify/functions/flow-alert
```

with parameters like:

- `flowID`: The unique identifier for the flow.
- `owner`: The flow owner’s address
  - for ICS20 hooks this is be derived from the sender if the message is not going back to the same chain
- `email`: The recipient's email address (optional).
- `unsubscribe`: (optional) If `true`, unsubscribes the email from future alerts.

For example:

```
https://portal.intento.zone/.netlify/functions/flow-alert?flowID=60300&unsubscribe=true&email=john.doe%40gmail.com
```

This handles subscription management and flow alerting. Whitelisting of your url is required for the live system—**message us to get access**.

### **Customizing the Alert Page (`/alert`)**

Integrators and their users can point towards the the alert page at:

```
https://portal.intento.zone/alert
```

You can customize the display using URL parameters:

- `flowID`: The flow to display (e.g., `flowID=8`).
- `owner`: The flow owner’s name or label.
- `imageUrl`: A custom image to show.
- `title`: Override the page title.
- `description`: Add a custom description.
- `theme`: Choose between `light` or `dark` mode.

Example:

```
https://portal.intento.zone/alert?flowID=8&theme=dark&imageUrl=https://www.svgrepo.com/show/303106/mcdonald-s-15-logo.svg
```

This gives you a flexible, user-facing alert page that fits into your app, dashboard, or site.

