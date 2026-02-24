# Bugs and Fixes
1. The UI especially the landing pages look off in some parts. Like the landing page has too big right and left margins. The UI looks like ameature work. Rest redesign it, make it look like a unicorn product, without changing content, eliminate extra styling that are not need, improve the styling and make it standout. Make it simple to use, minimal yet powerful. No one need to read documentation to learn how to use it.
2. When user is logged in, do not show buttons like login, signup, or try for free. Let them proceed to create new event/stream and not try. Try if for users who have not logged in. Make it easy for them to access their dashboards.
3. Remove fake testimonials and numbers in signup pages. If numbers must be used, pull it from the db or something. Make it credible.
4. On the dashboard, remove "Start a Live Room" since it is already handled in "create event" section as a an event_type.
5. Show users the events they've paid/subscribed for. Let them track the upcoming and concluded events.
6. Withdrawal not showing. Users needed to add their account and when trigger withdrawal, we use IntaSend payout API to send funds to those accounts. They can choose to withdraw to either Bank or Mobile account. Must enter OTP sent on email to process withdrawal request.
7. Email like sign up (email verification) are not being sent out.
8. In the event_type = Video+Audio, video is not shown or no way of activating in the UI. The idea was to use LiveKit to create interactive video. Use the current event_type=audio only UI, and just add video enabling button and make it work.
9. Remain anything commentary to streams. No longer a sports only platform
10. Come up with a single logo and make an SVG file of it. Use it everywhere and make things consistent.
11. Text and fonts should be readable. Simplicity and usability wins. Remember my unicorn comparison.
12. Users can mark events public or private. Public events should be published on a page labelled Discover Streams or similar. We can list upcoming and recently completed like 20 of them. Private events, users must share the links to be visible. Public events will be discoverable and can be searched.
13. Keep things DRY where possible. Do not break things that work. Improve where you can. Fix bugs highlighted above. Feel free to ask questions