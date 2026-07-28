Elysium project draft

## Status

This is a design-discussion writeup, not a spec. Everything under
"Architecture" and "Core mechanics" is a sketch that came out of talking
through the idea, not a commitment. Expect it to change once real code
exists and things get eyeballed on an actual terminal. Treat concrete
type names, package names, and library choices as starting points, not
as decisions to defend.

## Working agreement for this project — read this first

- This project is hand-coded only. All code is written by the user,
  by hand.
- Claude Code (or any coding agent) is for **discussion only**:
  talking through approaches, explaining concepts, sketching
  pseudocode or type shapes on request, reviewing code the user
  wrote and pastes in. It is not for writing, generating, editing, or
  refactoring code automatically, and not for creating or editing
  project files directly.
- If asked something like "implement X" or "add Y", the right
  response is to talk through how to approach it (and sketch
  pseudocode inline in the chat if that helps), not to produce a
  diff, create a file, or edit one.
- This holds for the whole project, not just the first session.

## Concept

A continuously-running terminal simulation of an aquatic ecosystem
built around a specific real species, Elysia chlorotica (the "leaf
sea slug"), rendered with a cbonsai-like
visual sensibility: deliberate color palette, negative space, one
area of real detail (the creatures) against a calm background.

It is not a fish-tank screensaver in the asciiquarium sense (generic
species, decorative only) and not a snapshot/idle-game tool that
fast-forwards between launches — it's a live, continuously rendered
process, closer in spirit to asciiquarium's format but built around
one real organism's actual (loosely adapted) life cycle and a
population running at once, rather than generic decorative sprites.

No player goals, no win state, no completion. The mechanism that
prevents completion is population turnover (individuals are born and
die on staggered timelines) rather than any single creature's life
story ending.

## Species and biology grounding

Real biology used as design inspiration, not something the sim needs
to model with scientific rigor:

- Elysia chlorotica is stenophagous — feeds on one algae species,
  Vaucheria litorea, using a radula to pierce cells and suck out
  contents, but does not digest the chloroplasts, which get
  incorporated into its digestive tubules.
- Color shifts from reddish-brown to green over roughly 48 hours as
  chloroplasts accumulate.
- "Permanent kleptoplasty" (the state where it no longer needs to
  feed) requires roughly a week of *continuous* feeding; shorter
  feeding bouts cause the plastids to be lost again ("transient"
  kleptoplasty) — this is the origin of the "feeding streak, not just
  a fill bar" mechanic below.
- Adults can then run mostly on photosynthesis for a long stretch
  (reported around 10 months in the real animal).
- Light is not simply "more is better": in related species, high
  irradiance can shorten kleptoplast retention (photo-oxidative
  stress) while low light extends it — this is the origin of the
  "must find a light band, not just chase brightness" adult mechanic.
- Real habitat is very shallow (salt marshes, tidal pools, 0–0.5m),
  not open ocean. The project deliberately takes artistic license and
  sets this in open water instead — a conscious departure, not an
  oversight.
- None of the above needs to stay accurate as the design evolves.
  It's the seed the mechanics grew from, not a constraint on the
  implementation.

## Setting

Open water column in a salt marsh or tidal pool(doesn't need to be specified), not an enclosed tank. Continuously rendered in
real time

## Core mechanics (draft)

### Population, not one organism

Many individuals exist at once, at staggered life stages. This is
what makes the sim endlessly re-watchable without a completion state
and without needing to be stared at like a video — at any glance,
some individuals are mid-story, so there's no single arc to miss by
looking away.

### Reproduction — the portal

Each adult carries its own cooldown timer. When it fires, a circular
portal effect (concentric rings or a rotating spiral) renders at
*that specific adult's current position* for a second or two; when
it completes, a new juvenile materializes there with freshly
randomized traits. Tying the portal to a specific parent (rather
than a fixed spawn point) was a deliberate choice for legibility —
you can watch one adult and eventually see it "give birth."

### Juvenile stage — foraging and competition

- Food is individual algae particles (not static patches), drifting
  on a noise-driven flow field rather than randomly.
- Each juvenile gets randomized traits at spawn: initial size,
  max speed, max acceleration, and growth rate (independent of
  initial size — a juvenile can start small and grow fast, or start
  big and grow slow).
- Simple seek-steering toward the nearest/best food particle.
  Competition is implicit: whoever's steering and traits get them
  there first "wins" the particle — there is no explicit contest
  logic.
- `currentSize = initialSize + growthRate * fullness`
- Fullness climbs via a continuous feeding streak (mirrors real
  permanent vs. transient kleptoplasty). Stopping early causes
  fullness to backslide rather than hold — a real failure mode with
  no explicit punishment, just an ambient story beat.

### Adult stage — seeking the light band

- Crossing the fullness threshold locks in size, shifts color fully
  to green, unfurls the parapodia, and flips behavior from
  food-seeking to light-band-seeking.
- The "light band" is a horizontal zone, visually obvious (not a
  subtle tint), whose vertical center drifts over time (sine plus
  noise) standing in for changing light conditions.
- Adults ease toward the band lazily rather than snapping to it —
  movement quality should read as "loafing," in contrast to a
  juvenile's visibly hustling steering, so life stage is legible from
  motion alone, without needing labels.
- Adults have a timer which when expires, they show an animation of calling a
  portal, and a portal randomly generates somewhere. it will have a randomly
  chosen colour from a given range/palette which the adult's calling animation and portal both match
  The adult will then de-materialise and in code be forgotten about

### Open questions — not yet resolved

- **Population density**: how many concurrent juveniles / adults /
  food particles actually reads as "a rich elysium" vs. "empty" or
  "cluttered" for a typical terminal size. This can't be settled on
  paper — it needs eyeballing once something renders.

## Visual / aesthetic direction

Anchor reference: cbonsai — but specifically its *aesthetic*
(deliberate palette, negative space, one area of real intricacy
against a calm background), not its snapshot-then-done format.

- 24-bit color throughout, not 256-color or default ANSI.
- A considered, consistent palette per zone (water gradient, algae,
  adult green, portal effect) rather than default terminal colors.
- Rendering resolution: Unicode Braille characters (2×4 sub-cell dot
  grid per character) for smooth spine/parapodia curves, since
  creature positions are floating-point, not grid-locked.
- Cheap ambient texture (light shimmer, faint particulate drift) layered
  in without competing with the creatures for attention — the same
  role cbonsai's pot detailing plays.

## Procedural animation approach (Rain-World-inspired)

- **Body**: a short chain of point masses (roughly 8–12), advanced
  with Verlet integration plus iterative distance-constraint
  relaxation — standard rope/soft-body physics.
- **Locomotion**: a traveling wave down the spine
  (`sin(phase - k*index)` applied as a lateral force, not a position
  snap), modeling the real muscular pedal wave gastropods use to
  crawl — chosen to be the actual mechanism, not a decorative wiggle.
- **Parapodia**: a separate, simpler system — a spring-damped angle
  toward a light-dependent target (open when basking, folded when
  foraging or in the dark), not part of the main constraint chain.
- **Color**: interpolated between juvenile brown/red and adult green
  as a function of fullness, in a perceptual color space rather than
  naive RGB, to avoid a muddy midpoint.

  Maybe to start with we can use a somewhat more static animation rather than
  the procedural animation.

## Architecture (draft)

Two-package split as the core boundary: `sim` (pure logic, no I/O,
owns all persistent-per-frame state) and `render` (draws a `World`
snapshot, never mutates it).

```
elysium/
  go.mod
  cmd/elysium/main.go       — bubbletea program setup, tick loop, resize wiring

  internal/sim/
    world.go                 — World{Slugs, Food, t}; World.Tick(dt)
    slug.go                  — Slug: stage, traits, spine, parapodia, fullness, portal timer
    food.go                  — FoodParticle: pos, size
    environment.go           — LightBand(t), FlowAt(pos, t) via noise
    steering.go               — shared seek-steering (juvenile→food, adult→band)
    physics.go                — Verlet + constraint relaxation for spine chains

  internal/render/
    canvas.go                — custom Braille dot-grid + per-cell color
    palette.go               — perceptual color ramps
    frame.go                  — draws background → band → food → creatures
    portal.go                 — the ring/spiral effect
```

Libraries under consideration, each for a specific reason:

- `charmbracelet/bubbletea` — terminal setup/teardown, resize
  handling, tick loop.
- `charmbracelet/harmonica` — spring-based animation for the
  parapodia angle and adult drift toward the light band (not for the
  spine chain, which needs real constraint physics, not a single
  spring).
- `lucasb-eyer/go-colorful` — perceptual-space color interpolation,
  to avoid muddy RGB blending.
- `ojrac/opensimplex-go` — noise for the food/current flow field.
- A **hand-written** Braille+color canvas rather than an existing
  drawille port, because the existing Go ports (e.g.
  `exrook/drawille-go`) track shape but not per-cell color, and this
  project needs both.
- Deliberately no ECS framework — a dozen-ish entities in slices is
  the right scale.

Suggested build order: physics + steering first, proven via printed
coordinates with no rendering at all. Then the Braille+color canvas,
tested on a static shape in isolation. Only then combine them. These
two pieces (the spine physics and the canvas) are the ones likely to
eat unplanned time; proving each alone first is what keeps this
tractable.

## Deliberately deferred / out of scope for now

The "inputs" side of the project was intentionally split off and
parked, not forgotten:

- Terminal focus events (bubbletea's `WithReportFocus` /
  `FocusMsg` / `BlurMsg`) as a Forest-app-style focus-timer signal —
  leaving the terminal unfocused would feed into the environment as
  a stress state.
- Piping stdin through the tool as a progress-bar replacement for
  long-running commands.
- System stats (CPU/mem/network) driving tank conditions as an
  ambient dashboard.

None of this should be worked on until the simulation itself exists
and looks right on its own.

Also undecided, but low-stakes: whether simulation state persists
across process restarts.

## Origin, for reference only

Originally scoped as a Go CLI project, completable in about a month,
chosen to be scalable and eventually shifted toward "playful" after
an initial round of more utilitarian systems-tool ideas (P2P file
transfer, an ELF/DWARF crash-analysis tool, a network diagnostic
CLI) didn't land. A terminal life sim was chosen over a poker-equity
trainer idea, generalized from a generic fish tank to this
species-specific concept, and originally estimated at about a day —
that estimate is stale; the design grew substantially through
discussion, which is fine, just noted here so it doesn't cause
confusion later.
