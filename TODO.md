# Limitations of current version

As a course project, this repository is limited in scope. To avoid spending
unbounded time on implementation details, the following compromises are made:

## Backend

### Security & Privacy

- No password is needed to log into the site.
  - Having password authentication does not mean the system is secure. For demo
    purpose, we decided that it's better not do it than do it wrong.
- No token is needed to access user photos, either uploaded or generated.
  - Adding access control for those files would complicate interaction with AI
    models hosted by Replicate.
- Files are stored using UNIX permission 0666 (parent directory 0777).
  - This does not follow the minimal permissions principle.

### Performance

- All photos are loaded when requesting user profile without pagination.
  - Because each demo session uses a new user account, this is not a problem.
    See the frontend section below.
- When starting a game, All photos are shuffled and excess photos are discarded.
  - Same as above - there would never be too many photos in the demo account.

### Persistence

- Games are not saved to DB. Thus, this version has no progress tracking.

### Maintainability

- The layout of all 4 levels are hard-coded.
- Authentication routine (cookie validation) is not abstracted.

Since we have implemented an extensible framework for this webapp, those issues
can be solved without fundamental change to the project structure, given proper
human resource and real need.

## Frontend

### Functionality

- No game progress tracking (same as backend)
- No profile page or album management.
  - Even though backend has the progress of initial AI processing after upload,
    current frontend can't display it, and the user still has to guess whether
    AI has finished generating altered images.
- Frontend does not currently check the user profile for previously uploaded
  photos. Instead, after logging in (or implicit registering), user is sent to
  the upload page invariably.
  - We had to carefully choose a new distinct user name for each demo session
    to avoid unexpected photos showing up.
