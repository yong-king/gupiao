# Frontend Navigation Usability Spec

## 1. Background

The operations console must be usable after login. The login screen must be visually centered, and each sidebar entry must switch to a meaningful view instead of being decorative.

## 2. Goals

- Center the registration and login panel in the available browser viewport.
- Make every sidebar navigation item clickable.
- Render a distinct content view for watchlists, holdings, alert rules, refresh jobs, alerts, reports, account monitoring, and settings.
- Keep all flows reminder-only; do not add trading actions.

## 3. Non-Goals

- Do not introduce automatic buy, sell, or order submission controls.
- Do not require PostgreSQL or Redis just to fix the frontend usability issue.
- Do not replace the current static MVP shell with a full Vue build in this plan.

## 4. Functional Scope

### Must Have

- Login state uses a dedicated centered layout.
- Authenticated state uses the app shell layout with sidebar and content.
- Sidebar buttons update the active view and expose an active visual state.
- Active view persists locally so refresh keeps the last opened view.
- Each view includes useful text or a safe MVP action.
- Registration/login stores `user_id` for later authenticated API calls.

### Should Have

- Watchlist, holdings, rules, refresh, and alerts pages should expose Demo actions against the existing API.
- Settings page should show where backend, Agent, database, Redis, and deployment configuration live.

## 5. Testing

- Frontend unit tests cover the navigation view copy model.
- Frontend typecheck must pass after code changes.
- Browser verification must confirm login layout and sidebar switching.

## 6. Acceptance Criteria

- Given the user is not logged in When the app renders Then the auth card is centered in a single-column layout.
- Given the user is logged in When clicking any sidebar item Then the main content changes to that item's view.
- Given the user runs Demo actions Then they use the stored authenticated `user_id`.
- Given the app renders Then no automatic trading action is visible.

## 7. Definition Of Done

- Spec and plan are updated.
- Relevant frontend tests pass.
- Browser smoke test passes.
