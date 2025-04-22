# TODO List

1. File: [./main.go](./main.go), Line: 52

    implement mongo db persistence implementation

2. File: [./GeneticAlgorithm/Selection.go](./GeneticAlgorithm/Selection.go), Line: 3

    tournament selection

3. File: [./GeneticAlgorithm/GeneticAlgorithm.go](./GeneticAlgorithm/GeneticAlgorithm.go), Line: 26

    implement the whole genetic algorithm function

4. File: [./GeneticAlgorithm/FitnessFunction_test.go](./GeneticAlgorithm/FitnessFunction_test.go), Line: 31

    complete this empty schedule fitness test

5. File: [./GeneticAlgorithm/FitnessFunction_test.go](./GeneticAlgorithm/FitnessFunction_test.go), Line: 110

    complete this generated schedule fitness test

6. File: [./GeneticAlgorithm/RandomMutation.go](./GeneticAlgorithm/RandomMutation.go), Line: 308

    debug there is an error here (currently this mutation function is not being used)

7. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 146

    implement distribution types  front compressed distribution : start  for now the this is a brute force implementation, yet it is still enough considering the current low numbers of courses, instructors, rooms and sections. but I believe this still can be optimize in the future if we wanted to.  or just create another implementation for inidividual generation.

8. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 158

    the map below was added for debugging purposes only, remove when the code becomes stable and optimized.

9. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 486

    [implement below] search for general rooms that are available (consult first)

10. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 488

    [implement below] search available lab room for lecture subjects (consult first)

11. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 650

    when generating solutions while the genetic algorithm is running, we should also generate an index file to be use for querying each sections in the generated university schedules, make the generated schedule and index global for access.

12. File: [./GeneticAlgorithm/FitnessFunction.go](./GeneticAlgorithm/FitnessFunction.go), Line: 121

    implement preference heat map comparison based fitness function

13. File: [./Routes/RoutesV2/ScheduleDelete.go](./Routes/RoutesV2/ScheduleDelete.go), Line: 54

    use department_id for authentication later on.

14. File: [./Routes/RoutesV1/SchedulePost.go](./Routes/RoutesV1/SchedulePost.go), Line: 266

    on genetic algorithm error - just use normal schedule generation result

15. File: [./Routes/RoutesV1/ScheduleDelete.go](./Routes/RoutesV1/ScheduleDelete.go), Line: 34

    use department_id for authentication later on.

16. File: [./Routes/RoutesV1/ScheduleDelete.go](./Routes/RoutesV1/ScheduleDelete.go), Line: 103

    use department_id for authentication later on.

