# TODO List

1. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 133

    implement distribution types  front compressed distribution : start  for now the this is a brute force implementation, yet it is still enough considering the current low numbers of courses, instructors, rooms and sections. but I believe this still can be optimize in the future if we wanted to.  or just create another implementation for inidividual generation.

2. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 145

    the map below was added for debugging purposes only, remove when the code becomes stable and optimized.

3. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 339

    [implement below] search for general rooms that are available (consult first)

4. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 341

    [implement below] search available lab room for lecture subjects (consult first)

5. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 454

    when generating solutions while the genetic algorithm is running, we should also generate an index file to be use for querying each sections in the generated university schedules, make the generated schedule and index global for access.

6. File: [./Schedule/UniTimeTables.go](./Schedule/UniTimeTables.go), Line: 262

    implement horizontal checks - incomplete code below

7. File: [./Routes/RoutesV1/GetClassSchedule.go](./Routes/RoutesV1/GetClassSchedule.go), Line: 39

    use department_id for authentication later on.

8. File: [./Routes/RoutesV1/GetClassSchedule.go](./Routes/RoutesV1/GetClassSchedule.go), Line: 176

    use department_id for authentication later on.

9. File: [./Routes/RoutesV1/GenerateSchedule.go](./Routes/RoutesV1/GenerateSchedule.go), Line: 35

    this is just for testing

