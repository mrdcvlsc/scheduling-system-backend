# TODO List

1. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 133

    implement distribution types  front compressed distribution : start  for now the this is a brute force implementation, yet it is still enough considering the current low numbers of courses, instructors, rooms and sections. but I believe this still can be optimize in the future if we wanted to.  or just create another implementation for inidividual generation.

2. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 145

    the map below was added for debugging purposes only, remove when the code becomes stable and optimized.

3. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 271

    debug later because even though this should speed up the algorithm, it seems that it is not happening often so it does not provide considerable amout of speed in performance.

4. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 289

    This is bad design, move the room search loop outside the instructor search loop in the futrue. even though the time complexity is now correct, people might get confused when they look at the code in the future since the room search loop is inside the instructor search loop, they might assume that the time complexity of this inner algorithm is O(I * R), even if this is actually just O(I + R),  I = total number of departamental instructor. R = total number of departamental rooms for a specific room type.

5. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 360

    [implement below] search for general rooms that are available (consult first)

6. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 362

    [implement below] search available lab room for lecture subjects (consult first)

7. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 456

    debug later, for some reasons there are some subjects not being assigned

8. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 477

    keep track of all the subjects that was assigned, find out which subject was not assigned and analyze it why is that happening.

9. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 494

    when generating solutions while the genetic algorithm is running, we should also generate an index file to be use for querying each sections in the generated university schedules, make the generated schedule and index global for access.

10. File: [./Schedule/UniTimeTables.go](./Schedule/UniTimeTables.go), Line: 262

    implement horizontal checks - incomplete code below

11. File: [./Routes/RoutesV1/GetClassSchedule.go](./Routes/RoutesV1/GetClassSchedule.go), Line: 23

    use department_id for authentication later on.

12. File: [./Routes/RoutesV1/GenerateSchedule.go](./Routes/RoutesV1/GenerateSchedule.go), Line: 13

    generate schedule if there is no schedule

13. File: [./Routes/RoutesV1/GenerateSchedule.go](./Routes/RoutesV1/GenerateSchedule.go), Line: 16

    the generate function should be a go routine

14. File: [./Routes/RoutesV1/GenerateSchedule.go](./Routes/RoutesV1/GenerateSchedule.go), Line: 35

    this is just for testing

