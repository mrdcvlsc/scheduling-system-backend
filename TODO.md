# TODO List

1. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 248

    debug later because even though this should speed up the algorithm, it seems that it is not happening often so it does not provide considerable amout of speed in performance.

2. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 266

    This is bad design, move the room search loop outside the instructor search loop in the futrue. even though the time complexity is now correct, people might get confused when they look at the code in the future since the room search loop is inside the instructor search loop, they might assume that the time complexity of this inner algorithm is O(I * R), even if this is actually just O(I + R),  I = total number of departamental instructor. R = total number of departamental rooms for a specific room type.

3. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 337

    [implement below] search for general rooms that are available (consult first)

4. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 339

    [implement below] search available lab room for lecture subjects (consult first)

5. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 442

    when generating solutions while the genetic algorithm is running, we should also generate an index file to be use for querying each sections in the generated university schedules, make the generated schedule and index global for access.

6. File: [./Schedule/UniTimeTables.go](./Schedule/UniTimeTables.go), Line: 263

    implement horizontal checks - incomplete code below

