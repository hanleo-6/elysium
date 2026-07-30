Elysium project design

## Working agreement for this project - read this first

- This project is hand-coded only. All code is written by the user,
  by hand.
- Any coding agent are used for **discussion only**:
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
- Exception: DEVPLAN.md and this file (CONTEXT.md) are process/
  planning documents, not simulation code, and the agent may create
  and edit them directly on request (e.g. to track checkpoint status
  or log a design decision from discussion). This does not extend to
  anything under `internal/`, `cmd/`, or any `.go` file - those
  remain hand-coded only.
- The user is new to Go but not to programming. Be ready to explain syntax or
  go-specific details, but assume there is no need to explain general
  programming concepts

See DEVPLAN.md for the actively-tracked implementation checkpoint
list and status - that file is the source of truth for build order
and progress; this file stays focused on design.

## Concept

A continuously-running terminal simulation of an aquatic ecosystem
built around a specific real species, Elysia chlorotica (the "leaf
sea slug"), rendered with a cbonsai-like
visuals: deliberate color palette, negative space
Takes some inspiration from fish-tank screensaver asciiquarium

## Species and biology grounding

Real biology used as design inspiration (but the sim doesn't need to model this
exactly):

- Elysia chlorotica is stenophagous - feeds on one algae species,
  Vaucheria litorea, using a radula to pierce cells and suck out
  contents, but does not digest the chloroplasts, which get
  incorporated into its digestive tubules.
- Color shifts from reddish-brown to green over roughly 48 hours as
  chloroplasts accumulate.
- "Permanent kleptoplasty" (the state where it no longer needs to
  feed) requires roughly a week of *continuous* feeding; shorter
  feeding bouts cause the plastids to be lost again ("transient"
  kleptoplasty).
- Adults can then run mostly on photosynthesis for a long stretch
  (reported around 10 months in the real animal).
- Light is not simply "more is better": in related species, high
  irradiance can shorten kleptoplast retention (photo-oxidative
  stress) while low light extends it - this is the origin of the
  "must find a light band" adult mechanic.
- Real habitat is very shallow (salt marshes, tidal pools, 0–0.5m),
  not open ocean.

## Setting

Water column in a salt marsh.

## Core mechanics (draft)

### Population
Many individuals exist at once, at staggered life stages. This is
what makes the sim endlessly re-watchable without a completion state
and without needing to be stared at like a video - at any glance,
some individuals are mid-story, so there's no single arc to miss by
looking away.

### Reproduction - the portal

Each adult carries its own timer. When it runs out, a circular
portal effect (rotating spiral) renders at
*that specific adult's current position* for a second or two; then,
a juvenile materializes there with freshly randomized traits.

### Juvenile stage - foraging and competition

- Food is individual algae particles, drifting
  on a noise-driven flow field.
- Each juvenile gets randomized traits at spawn: initial size,
  max speed, max acceleration, and growth rate
- Simple seek-steering toward the nearest/best food particle.
  Competition is implicit: whoever's steering and traits get them
  there first "wins" the particle
- `currentSize = initialSize + growthRate * fullness`
- Fullness climbs via a continuous feeding streak (mirrors real
  permanent vs. transient kleptoplasty). Stopping early causes
  fullness to decay.

### Adult stage - seeking the light band

- Crossing the fullness threshold locks in size, shifts color fully
  to green, unfurls the parapodia, and flips behavior from
  food-seeking to light-band-seeking.
- The "light band" is a horizontal zone, visually obvious,
  whose vertical center drifts over time (sine plus
  noise) standing in for changing light conditions.
- Adults loaf toward the band lazily, in contrast to a
  juvenile's hustle steering

## Visual / aesthetic direction

Reference 1: cbonsai
### Minimal and vibrant
Specifically its *aesthetic* (deliberate palette, negative space, one area of real intricacy
against a calm background), not its snapshot-then-done format.

- 24-bit color throughout, not 256-color or default ANSI.
- A considered, consistent palette per zone (water gradient, algae,
  adult green, portal effect) rather than default terminal colors.
- Rendering resolution: Unicode Braille characters (2×4 sub-cell dot
  grid per character) for smooth spine/parapodia curves
- Ambient texture (light shimmer, faint particulate drift) layered
  in without competing with the creatures for attention

Reference 2: Rain World
### Procedural animation

- **Body**: a short chain of point masses (roughly 8–12), standard rope/soft-body physics.
  Using Verlet integration plus iterative distance-constraint relaxation
- **Locomotion**: a traveling wave down the spine
  (`sin(phase - k*index)` applied as a lateral force, not a position
  snap), modeling the real muscular pedal wave gastropods use to
  crawl.
- **Parapodia**: a separate, simpler system - a spring-damped angle
  toward a light-dependent target (open when basking, folded when
  foraging or in the dark), not part of the main constraint chain.
- **Color**: interpolated between juvenile brown/red and adult green
  as a function of fullness, in perceptual color space.

## Architecture (working draft)

```
elysium/
  go.mod
  cmd/elysium/main.go       - bubbletea program setup, tick loop, resize wiring

  internal/sim/
    saltmarsh.go                 - World{Slugs, Food, t}; World.Tick(dt)
    seaslug.go                  - Slug: stage, traits, spine, parapodia, fullness, portal timer
    algae.go                  - FoodParticle: pos, size
    environment.go           - LightBand(t), FlowAt(pos, t) via noise
    seek.go               - juvenile -> food, adult -> light band

  internal/physics/
    physics.go                - Verlet + constraint relaxation for spine chains

  internal/render/
    canvas.go                - custom Braille dot-grid + cell colouring
    palette.go               - colour definitions
    frame.go                  - draws each frame: background → band → food → creatures
    portal.go                 - spiral portal
    elysia_chlorotica.go      - elysia_chlorotica procedural rendering
    sprites.go                - sprites for static elements (algae, extra environment objects?)
    shapes.go                 - primitive shapes

```

Libraries under consideration:

- `charmbracelet/bubbletea` - terminal setup/teardown, resize
  handling, tick loop.
- `charmbracelet/harmonica` - spring-based animation for the
  parapodia angle and adult drift toward the light band (not for the
  spine chain
- `lucasb-eyer/go-colorful` - perceptual-space color interpolation,
- `ojrac/opensimplex-go` - noise for the food/current flow field.

## Out of scope for now

The "inputs" side of the project:

- Terminal focus events (bubbletea's `WithReportFocus` /
  `FocusMsg` / `BlurMsg`) as a Forest-app-style focus-timer signal
- Piping stdin through the tool as a progress-bar replacement for
  long-running commands.
- System stats (CPU/mem/network) driving tank conditions as an
  ambient dashboard.

