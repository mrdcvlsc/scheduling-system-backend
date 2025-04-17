# TODO List

1. File: [./main.go](./main.go), Line: 52

    implement mongo db persistence implementation

2. File: [./GeneticAlgorithm/Selection.go](./GeneticAlgorithm/Selection.go), Line: 3

    tournament selection

3. File: [./GeneticAlgorithm/GeneticAlgorithm.go](./GeneticAlgorithm/GeneticAlgorithm.go), Line: 15

    implement the whole genetic algorithm function

4. File: [./GeneticAlgorithm/GeneticAlgorithm.go](./GeneticAlgorithm/GeneticAlgorithm.go), Line: 113

    crossover some individuals and add them to the CURRENT population

5. File: [./GeneticAlgorithm/GeneticAlgorithm.go](./GeneticAlgorithm/GeneticAlgorithm.go), Line: 115

    apply random mutations to some of the CURRENT individuals in the population

6. File: [./GeneticAlgorithm/GeneticAlgorithm.go](./GeneticAlgorithm/GeneticAlgorithm.go), Line: 117

    after random mutations and crossover, some individuals might be incomplete so we need to re-encode them again to complete their schedules.

7. File: [./GeneticAlgorithm/GeneticAlgorithm.go](./GeneticAlgorithm/GeneticAlgorithm.go), Line: 120

    apply tournament selection to be saved to the next generation

8. File: [./GeneticAlgorithm/GeneticAlgorithm.go](./GeneticAlgorithm/GeneticAlgorithm.go), Line: 133

    sort individuals base on fitness

9. File: [./GeneticAlgorithm/Crossover.go](./GeneticAlgorithm/Crossover.go), Line: 3

    implement a crossover function

10. File: [./GeneticAlgorithm/Crossover.go](./GeneticAlgorithm/Crossover.go), Line: 5

    1. just crossover the two parents

11. File: [./GeneticAlgorithm/Crossover.go](./GeneticAlgorithm/Crossover.go), Line: 7

    2. check for overlapping vertical errors and horizontal errors

12. File: [./GeneticAlgorithm/Crossover.go](./GeneticAlgorithm/Crossover.go), Line: 9

    3. delete the other genes or data that overlaps (pick a random gene/data that will remained)

13. File: [./GeneticAlgorithm/Crossover.go](./GeneticAlgorithm/Crossover.go), Line: 11

    4. re-encode to complete the offspring's genome

14. File: [./GeneticAlgorithm/FitnessFunction_test.go](./GeneticAlgorithm/FitnessFunction_test.go), Line: 24

    complete this empty schedule fitness test

15. File: [./GeneticAlgorithm/FitnessFunction_test.go](./GeneticAlgorithm/FitnessFunction_test.go), Line: 103

    complete this generated schedule fitness test

16. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 146

    implement distribution types  front compressed distribution : start  for now the this is a brute force implementation, yet it is still enough considering the current low numbers of courses, instructors, rooms and sections. but I believe this still can be optimize in the future if we wanted to.  or just create another implementation for inidividual generation.

17. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 158

    the map below was added for debugging purposes only, remove when the code becomes stable and optimized.

18. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 455

    [implement below] search for general rooms that are available (consult first)

19. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 457

    [implement below] search available lab room for lecture subjects (consult first)

20. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 588

    when generating solutions while the genetic algorithm is running, we should also generate an index file to be use for querying each sections in the generated university schedules, make the generated schedule and index global for access.

21. File: [./GeneticAlgorithm/FitnessFunction.go](./GeneticAlgorithm/FitnessFunction.go), Line: 92

    implement preference heat map comparison based fitness function

22. File: [./Routes/RoutesV2/ScheduleDelete.go](./Routes/RoutesV2/ScheduleDelete.go), Line: 54

    use department_id for authentication later on.

23. File: [./Routes/RoutesV1/ScheduleDelete.go](./Routes/RoutesV1/ScheduleDelete.go), Line: 34

    use department_id for authentication later on.

24. File: [./Routes/RoutesV1/ScheduleDelete.go](./Routes/RoutesV1/ScheduleDelete.go), Line: 103

    use department_id for authentication later on.

