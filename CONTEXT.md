Elysium project design

> **Agents:** this is a design reference, not your instructions. The
> working agreement (hand-coded only, discussion-only agents, which
> docs you may edit) lives in `AGENTS.md` - follow that first.

This file is the design document: concept, biology grounding,
mechanics, aesthetic direction, and the current file architecture.
It's written mainly as context for coding agents, but stays readable
for humans.

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
- Adults seek toward the band, then loaf about in the zone lazily
  , in contrast to a juvenile's hustle steering.
- Light band implemented as an angled cone from a surface origin. Its
  angle sways via a sum of different phase sines plus mean-reverting
  (Ornstein–Uhlenbeck) drift, and its width breathes.

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

## Architecture (may be outdated)

Module: `github.com/hanleo-6/elysium` (Go 1.26). Dependency direction is
`cmd → render → sim → physics`; `physics` imports nothing internal.

```
elysium/
  go.mod / go.sum
  .gitignore                    - ignores /elysium binary, debug.log, DEVPLAN.md
  CONTEXT.md                    - design (this file)
  DEVPLAN.md                    - checkpoint list + status (gitignored)

  cmd/elysium/
    main.go                     - bubbletea v2 model: owns canvas, Environment,
                                  []*ElysiaChlorotica, []*Alga. Init/config
                                  (numSlugs, initialAlgae), fixed-dt tick
                                  (33ms / dt=0.033), resize → new Canvas,
                                  RegenerateAlgae (respawns eaten algae at random
                                  positions), View() draw order:
                                  Clear → algae → slugs → light band → Render.
                                  Logs to debug.log via tea.LogToFile.

  internal/physics/
    physics.go                  - Vec2 (Add/Sub/Scale/Length), Point{Pos, PrevPos},
                                  Constraint{P1, P2, RestLength}
    skeleton.go                 - Skeleton: NewSpine (straight point chain),
                                  Step (Verlet on head with accel, free Verlet on
                                  body, 5 constraint-relaxation iterations),
                                  Seek (Reynolds seek steering, force-clamped)

  internal/sim/
    environment.go              - Environment{grid, Saltmarsh, Time}: analytic
                                  sine "potential" → curl-noise flow (sampleCurl),
                                  flowAt = curl + wake; Tick advances wake decay,
                                  drifts algae, injects slug head velocity into
                                  wake grid, moves light band, advances Time
    wakegrid.go                 - wakeGrid: coarse vector grid of slug wakes
                                  (injectForce, bilinear sampleWake, decay tick)
    saltmarsh.go                - Saltmarsh: light band state (origin, InitTheta,
                                  SpreadTheta); MoveLightband(t, dt) = sum of 3
                                  sines + Ornstein–Uhlenbeck drift on the centre
                                  angle, sine-modulated spread
    seaslug.go                  - ElysiaChlorotica{Skelly, Fullness, BodyProfile,
                                  Stage}: Tick → SeekFood (nearest uneaten alga,
                                  eat within nomRange, Fullness += NomFactor),
                                  WingsOpenFactor. Stage enum exists, not yet
                                  driven.
    algae.go                    - Alga{Pos, Radius, Eaten}; NomFactor from radius

  internal/render/
    canvas.go                   - Braille canvas: 2×4 dots per cell, one fg + one
                                  bg colour per cell; SetDot/SetDotWithColor/
                                  SetBackground, FillBackground (water gradient,
                                  currently unused), Clear, Render (24-bit ANSI)
    palette.go                  - palette constants + SlugColor(fullness),
                                  AlgaColor(), WaterColor(depth) via go-colorful
    shapes.go                   - DrawCircle / DrawCircleFill (midpoint circle),
                                  DrawCone (fan of lines, writes background),
                                  DrawLine (Bresenham, pluggable set-point fn)
    lightband.go                - DrawLightband: degrees → radians, DrawCone
    alga.go                     - DrawAlga: filled circle
    elysia_chlorotica.go        - DrawECSlug: body circles from BodyProfile,
                                  parapodia lines offset along spine normals,
                                  rhinopores at head
    spine.go                    - DrawSpine: debug view of skeleton constraints
    sprites.go                  - AlgaeSprite (unused)
    frame.go                    - stub (intended home for main.go's draw loop)
    wormhole.go                 - stub (intended portal effect)
```

Planned, not yet created (from the original draft):

- `internal/sim/seek.go` - adult → light-band seeking (juvenile → food

Libraries:

- In use: `charm.land/bubbletea/v2` (terminal, resize, tick loop),
  `lucasb-eyer/go-colorful` (Lab/HCL colour blending).
- Still under consideration: `charmbracelet/harmonica` (spring for
  parapodia angle / adult drift).
- Dropped for now: `ojrac/opensimplex-go` - flow field uses an analytic
  sum-of-sines potential with curl instead.

## Out of scope for now

The "inputs" side of the project:

- Piping stdin into it as a progress-bar replacement for
  long-running commands.
- System stats (CPU/mem/network) driving tank conditions as an
  ambient dashboard.

