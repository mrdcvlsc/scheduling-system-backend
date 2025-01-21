# TODO List

1. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 208

    if this is the last time slot and there is still no instructor available, throw an error saying there is not enough instructors, true error handling not panic.

2. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 283

    iterate over all of sorted rooms to find which one is available, if there is no room available for the current time slot, continue to the next iteration of the loop.

3. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 287

    if this is the last time slot and there is still no room available, throw an error saying there is not enough room.

4. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 290

    after all of the checks and searches, if there are instructors and rooms that are available for a specific subject and time slot, then assign the selected subject, instructor and room to the time slot.

5. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 333

    during the first schedule generation throw error if there is not enough instructors

6. File: [./GeneticAlgorithm/GeneratePopulations.go](./GeneticAlgorithm/GeneratePopulations.go), Line: 335

    during the first schedule generation throw error if there is not enough rooms

7. File: [./Resources/Rooms/Rooms.go](./Resources/Rooms/Rooms.go), Line: 8

    (we only need 4 bits to store up to 15 hours) - Test with array dim : [6][24]

8. File: [./Resources/Schedule/UniTimeTables.go](./Resources/Schedule/UniTimeTables.go), Line: 41

    enable the code below after room assigning is implemented.

9. File: [./Resources/Schedule/UniTimeTables.go](./Resources/Schedule/UniTimeTables.go), Line: 51

    enable the code below after room assigning is implemented.

10. File: [./Resources/Schedule/UniTimeTables.go](./Resources/Schedule/UniTimeTables.go), Line: 74

    enable the code below after room assigning is implemented.

11. File: [./Resources/Schedule/UniTimeTables.go](./Resources/Schedule/UniTimeTables.go), Line: 104

    enable the code below after room assigning is implemented.

12. File: [./Resources/Schedule/UniTimeTables.go](./Resources/Schedule/UniTimeTables.go), Line: 136

    enable the code below after room assigning is implemented.

13. File: [./Resources/Schedule/UniTimeTables.go](./Resources/Schedule/UniTimeTables.go), Line: 172

    implement horizontal checks

